import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: string | Date): string {
  return new Date(date).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

export function formatDateTime(date: string | Date): string {
  return new Date(date).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function getInitials(firstName: string, lastName: string): string {
  return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase()
}

export function getPriorityColor(priority: string): string {
  const colors: Record<string, string> = {
    urgent: 'text-red-600 bg-red-50',
    high: 'text-orange-600 bg-orange-50',
    medium: 'text-yellow-600 bg-yellow-50',
    low: 'text-blue-600 bg-blue-50',
    none: 'text-gray-600 bg-gray-50',
  }
  return colors[priority] || colors.none
}

export function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    'backlog': 'bg-gray-200 text-gray-800',
    'to do': 'bg-blue-200 text-blue-800',
    'in progress': 'bg-yellow-200 text-yellow-800',
    'in review': 'bg-purple-200 text-purple-800',
    'done': 'bg-green-200 text-green-800',
  }
  return colors[status.toLowerCase()] || 'bg-gray-200 text-gray-800'
}
