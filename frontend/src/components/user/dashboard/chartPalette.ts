// Chart colors are separate from the monochrome interface theme.
// Keep Token categories consistent between the composition and trend charts.
export function chartPalette(dark: boolean) {
  return dark
    ? { total: '#a5b4fc', parts: ['#60a5fa', '#fbbf24', '#f472b6', '#a78bfa'], cache: '#fb7185', cost: '#fbbf24', requests: '#60a5fa' }
    : { total: '#4f46e5', parts: ['#2563eb', '#b45309', '#db2777', '#7c3aed'], cache: '#e11d48', cost: '#b45309', requests: '#2563eb' }
}

export function distributionColorIndex(id: string) {
  let hash = 0
  for (const char of id) hash = (Math.imul(hash, 31) + char.charCodeAt(0)) | 0
  return (hash >>> 0) % 6
}
