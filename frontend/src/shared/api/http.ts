// Модуль общего HTTP-хелпера для fetch-обёрток бэкенда.
// Базовый URL — относительный путь '/'; в dev запросы проксируются
// Vite-сервером на бэкенд (см. server.proxy в vite.config.ts).

/**
 * Выполняет fetch-запрос и разбирает JSON-ответ.
 * При не-2xx статусе выбрасывает ошибку с сообщением от бэкенда ({"error":"..."}).
 * Ответы без тела (например, 204) возвращают undefined.
 */
export async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init)

  if (!res.ok) {
    let message = `HTTP ${res.status} ${res.statusText}`
    try {
      const body = (await res.json()) as { error?: string }
      if (body.error) {
        message = body.error
      }
    } catch {
      // тело не JSON — оставляем дефолтное сообщение
    }
    throw new Error(message)
  }

  if (res.status === 204) {
    return undefined as T
  }

  return res.json() as Promise<T>
}