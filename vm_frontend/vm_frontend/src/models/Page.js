export class Page {
  constructor({ id, name, description = '', component, parent = null }) {
    this.id = id
    this.name = name
    this.description = description
    this.component = component
    this.parent = parent
    this.type = 'page'
  }

  get path() {
    const parentPath = this.parent ? this.parent.path : ''
    return `${parentPath}/${this.id}`
  }

  get routeConfig() {
    return {
      path: this.path,
      name: this.path.replace(/\//g, '_'),
      component: this.component
    }
  }
}

export class Folder {
  constructor({ id, name, description = '', children = [], parent = null }) {
    this.id = id
    this.name = name
    this.description = description
    this.children = children
    this.parent = parent
    this.type = 'folder'
    
    // Установить parent для всех детей
    this.children.forEach(child => {
      child.parent = this
    })
  }

  get path() {
    const parentPath = this.parent ? this.parent.path : ''
    return `${parentPath}/${this.id}`
  }

  getAllRoutes() {
    const routes = []
    
    for (const child of this.children) {
      if (child.type === 'page') {
        routes.push(child.routeConfig)
      } else if (child.type === 'folder') {
        routes.push(...child.getAllRoutes())
      }
    }
    
    return routes
  }
}
