export function syncThemeColor(isDark: boolean): void {
  const themeColorMeta = document.querySelector<HTMLMetaElement>('meta[name="theme-color"]')
  if (themeColorMeta) {
    themeColorMeta.content = isDark ? '#090a0c' : '#f5f6fb'
  }
}
