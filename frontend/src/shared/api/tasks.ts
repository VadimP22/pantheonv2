// Модуль fetch-обёрток для HTTP API задач.
import type { CreateTaskInputDTO, DepositTaskInputDTO, TaskDTO, TaskPageDTO } from '../generated/dto.ts'
import { request } from './http.ts'

/** Параметры пагинации списка задач. */
export interface ListTasksParams {
  /** Номер страницы (1-based). */
  page?: number
  /** Размер страницы (1..100). */
  pageSize?: number
}

/**
 * Возвращает задачу по идентификатору.
 * GET /api/tasks/{id}
 */
export async function getTask(id: number): Promise<TaskDTO> {
  return request<TaskDTO>(`/api/tasks/${id}`)
}

/**
 * Вносит ресурс в задачу.
 * POST /api/tasks/{id}/deposit
 */
export async function depositTask(id: number, input: DepositTaskInputDTO): Promise<TaskDTO> {
  return request<TaskDTO>(`/api/tasks/${id}/deposit`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
}

/**
 * Создаёт новую задачу.
 * POST /api/tasks
 */
export async function createTask(input: CreateTaskInputDTO): Promise<TaskDTO> {
  return request<TaskDTO>('/api/tasks', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
}

/**
 * Возвращает страницу задач.
 * GET /api/tasks
 */
export async function listTasks(params: ListTasksParams = {}): Promise<TaskPageDTO> {
  const { page = 1, pageSize = 6 } = params
  const qs = new URLSearchParams({
    page: String(page),
    page_size: String(pageSize),
  })
  return request<TaskPageDTO>(`/api/tasks?${qs.toString()}`)
}

/**
 * Удаляет задачу по идентификатору.
 * DELETE /api/tasks/{id}
 */
export async function deleteTask(id: number): Promise<void> {
  await request<void>(`/api/tasks/${id}`, { method: 'DELETE' })
}
