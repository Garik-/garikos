export function getColor(name: string, shade = null) {
  const varName = shade ? `--tblr-${name}-${shade}` : `--tblr-${name}`

  return getComputedStyle(document.documentElement).getPropertyValue(varName).trim()
}
