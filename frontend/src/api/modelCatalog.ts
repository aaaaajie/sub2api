export class ModelCatalogError extends Error {
  constructor(
    public kind: 'http' | 'network' | 'invalid',
    public status?: number,
    public code = '',
    public detail = ''
  ) { super(kind) }
}

export async function fetchModelCatalog(baseUrl: string, apiKey: string, signal?: AbortSignal): Promise<string[]> {
  const url = new URL(baseUrl || window.location.origin)
  const path = url.pathname.replace(/\/+$/, '')
  url.pathname = `${path.endsWith('/v1') ? path : `${path}/v1`}/models`
  url.search = ''
  url.hash = ''
  let response: Response
  try {
    response = await fetch(url.toString(), {
      headers: { Authorization: `Bearer ${apiKey}`, Accept: 'application/json' },
      credentials: 'omit', signal
    })
  } catch (error) {
    if (signal?.aborted) throw error
    throw new ModelCatalogError('network')
  }
  let payload: unknown
  try { payload = await response.json() } catch {
    throw new ModelCatalogError(response.ok ? 'invalid' : 'http', response.status)
  }
  const object = payload && typeof payload === 'object' ? payload as Record<string, unknown> : {}
  if (!response.ok) {
    const nested = object.error && typeof object.error === 'object' ? object.error as Record<string, unknown> : object
    // Never echo HTML or credentials, even if a proxy accidentally includes them.
    const safe = (value: unknown) => typeof value === 'string' && !/[<>]/.test(value)
      ? value.split(apiKey).join('[redacted]').replace(/sk-[\w-]+/g, '[redacted]').replace(/Bearer\s+\S+/gi, 'Bearer [redacted]').slice(0, 300) : ''
    throw new ModelCatalogError('http', response.status, safe(nested.code ?? nested.type), safe(nested.message))
  }
  const records = Array.isArray(object.data) ? object.data : Array.isArray(object.models) ? object.models : null
  if (!records) throw new ModelCatalogError('invalid')
  const models = records.map(record => {
    if (typeof record === 'string') return record.trim()
    if (!record || typeof record !== 'object') return ''
    return typeof record.id === 'string' ? record.id.trim() : ''
  }).filter(Boolean)
  if (records.length && !models.length) throw new ModelCatalogError('invalid')
  return [...new Set(models)]
}
