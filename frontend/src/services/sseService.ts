export function createSSEConnection(url: string, onMessage: (data: unknown) => void) {
  const eventSource = new EventSource(url)

  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data)
    onMessage(data)
  }

  eventSource.onerror = (error) => {
    console.error('SSE Error:', error)
    eventSource.close()
  }

  return eventSource
}
