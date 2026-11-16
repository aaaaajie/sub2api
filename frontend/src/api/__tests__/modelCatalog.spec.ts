import { afterEach, describe, expect, it, vi } from 'vitest'
import { fetchModelCatalog } from '../modelCatalog'
afterEach(() => vi.unstubAllGlobals())
const respond = (body: unknown, status = 200) => vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify(body), { status })))
describe('CCS model catalog', () => {
  it('reads models without generating and sends the key only in its header', async () => {
    respond({ data: [{ id: 'ds' }, { id: 'ds' }, { id: 'gpt' }] })
    expect(await fetchModelCatalog('https://example.com/api/v1/', 'sk-test')).toEqual(['ds', 'gpt'])
    expect(fetch).toHaveBeenCalledWith('https://example.com/api/v1/models', expect.objectContaining({ credentials: 'omit', headers: { Authorization: 'Bearer sk-test', Accept: 'application/json' } }))
  })
  it.each([
    [{ code: 'API_KEY_EXPIRED', message: 'Key expired' }, 'API_KEY_EXPIRED'],
    [{ error: { code: 'GROUP_NOT_ALLOWED', message: 'No permission' } }, 'GROUP_NOT_ALLOWED']
  ])('preserves structured HTTP reasons', async (body, code) => {
    respond(body, 403)
    await expect(fetchModelCatalog('https://example.com', 'sk-test')).rejects.toMatchObject({ kind: 'http', status: 403, code })
  })
  it('does not echo credentials', async () => {
    respond({ message: 'Bearer sk-test sk-other', code: 'DENIED' }, 401)
    await expect(fetchModelCatalog('https://example.com', 'sk-test')).rejects.toMatchObject({ detail: 'Bearer [redacted] [redacted]' })
  })
  it('does not render HTML errors', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('<html>bad gateway</html>', { status: 502 })))
    await expect(fetchModelCatalog('https://example.com', 'sk-test')).rejects.toMatchObject({ kind: 'http', status: 502, detail: '' })
  })
  it('distinguishes empty catalog from invalid payload', async () => {
    respond({ data: [] }); expect(await fetchModelCatalog('https://example.com', 'sk-test')).toEqual([])
    respond({ unexpected: true }); await expect(fetchModelCatalog('https://example.com', 'sk-test')).rejects.toMatchObject({ kind: 'invalid' })
  })
  it('distinguishes network failures and supports retry', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))
    await expect(fetchModelCatalog('https://example.com', 'sk-test')).rejects.toMatchObject({ kind: 'network' })
    respond({ data: [{ id: 'ds' }] }); expect(await fetchModelCatalog('https://example.com', 'sk-test')).toEqual(['ds'])
  })
  it('passes cancellation through', async () => {
    const controller = new AbortController(); controller.abort()
    const error = new DOMException('Aborted', 'AbortError')
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(error))
    await expect(fetchModelCatalog('https://example.com', 'sk-test', controller.signal)).rejects.toBe(error)
  })
})
