export const formatter = new Intl.NumberFormat('ru-RU', {
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
})

export const dateToLocaleString = (date: Date) =>
  date.toLocaleString('ru-RU', {
    dateStyle: 'short',
    timeStyle: 'short',
  })

export function formatBytes(bytes: number) {
  const units: ReadonlyArray<{ unit: string; value: number }> = [
    { unit: 'байт', value: 1 },
    { unit: 'Кбайт', value: 1024 },
    { unit: 'Мбайт', value: 1024 * 1024 },
    { unit: 'Гбайт', value: 1024 * 1024 * 1024 },
    { unit: 'Tбайт', value: 1024 * 1024 * 1024 * 1024 },
  ]

  let selectedUnit = units[0] ?? { unit: 'байт', value: 1 }
  for (let i = units.length - 1; i >= 0; i--) {
    const unit = units[i]
    if (!unit) continue
    if (bytes >= unit.value) {
      selectedUnit = unit
      break
    }
  }

  const value = bytes / selectedUnit.value
  const formattedValue = formatter.format(value)

  return formattedValue + ' ' + selectedUnit.unit
}
