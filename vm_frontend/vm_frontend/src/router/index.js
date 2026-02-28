import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/components/Home.vue'
import FolderView from '@/components/FolderView.vue'
import NotFound from '@/components/NotFound.vue'
import { structure, getAllRoutes } from '@/config/pages.js'

const routes = [
  { path: '/', name: 'Home', component: Home },
  // Динамические роуты для всех страниц из структуры
  ...getAllRoutes(structure),
  // Роуты для папок (должны быть после страниц)
  { path: '/:folder+', name: 'Folder', component: FolderView },
  // 404 для всех несуществующих страниц
  { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router