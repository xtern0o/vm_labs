import { Page, Folder } from '@/models/Page.js'

// структура с папками и страницами
export const structure = [
  new Folder({
    id: 'vm_labs',
    name: 'Вычислительная математика',
    description: 'Лабораторные работы по ВМ',
    children: [
      new Page({
        id: 'lab1',
        name: 'Лабораторная работа 1',
        description: 'Решение СЛАУ методом Гаусса-Зейделя',
        component: () => import('@/components/labs/Lab1.vue')
      }),
      new Page({
        id: 'lab2',
        name: 'Лабораторная работа 2',
        description: 'Решение нелинейных уравнений и систем',
        component: () => import('@/components/labs/Lab2.vue')
      }),
    ]
  }),
]

// сборка всех роутов из структуры
export function getAllRoutes(items) {
  const routes = []
  
  for (const item of items) {
    if (item.type === 'page') {
      routes.push(item.routeConfig)
    } else if (item.type === 'folder') {
      routes.push(...item.getAllRoutes())
    }
  }
  
  return routes
}
