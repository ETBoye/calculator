

export async function fetchWithErrorHandle<T>(
    input: string | URL | globalThis.Request,
    init?: RequestInit,
  ): Promise<T> {
    try {
      const response = await fetch(input, init)
  
      if (!response.ok) {
        alert(`Received response status: ${response.status} with body ${JSON.stringify(await response.json())}`)
        throw new Error(`Response not ok, status: ${response.status}`)
      }
  
      return await response.json()
  
    } catch(e) {
      alert("Received error from backend!");
      throw e;
    }
}

export function sessionIdIsValid(sessionId: string): boolean {
    return /^[a-zA-Z0-9\-]*$/.test(sessionId) && sessionId.length > 1 && sessionId.length < 100
}