import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import type { TaskPriority } from "@/types"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getPriorityColor(priority: TaskPriority): string {
  const colors = {
    urgent: 'text-red-600 border-red-600',
    high: 'text-orange-600 border-orange-600',
    medium: 'text-yellow-600 border-yellow-600',
    low: 'text-blue-600 border-blue-600',
    none: 'text-gray-600 border-gray-600',
  }
  return colors[priority] || colors.none
}

export function getStatusColor(status: string): string {
  const lowercaseStatus = status.toLowerCase()
  if (lowercaseStatus.includes('done') || lowercaseStatus.includes('completed')) {
    return 'bg-green-100 text-green-800'
  }
  if (lowercaseStatus.includes('progress') || lowercaseStatus.includes('doing')) {
    return 'bg-blue-100 text-blue-800'
  }
  if (lowercaseStatus.includes('review')) {
    return 'bg-purple-100 text-purple-800'
  }
  if (lowercaseStatus.includes('todo') || lowercaseStatus.includes('backlog')) {
    return 'bg-gray-100 text-gray-800'
  }
  return 'bg-gray-100 text-gray-800'
}

export function getInitials(firstName?: string, lastName?: string): string {
  if (!firstName && !lastName) return '??'
  const first = firstName?.charAt(0).toUpperCase() || ''
  const last = lastName?.charAt(0).toUpperCase() || ''
  return first + last || '?'
}

export function formatDate(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

export function formatDateTime(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
