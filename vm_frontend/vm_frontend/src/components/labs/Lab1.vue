<script setup>
import { ref } from 'vue'
import api from '@/services/api.js'
import Alert from '@/components/Alert.vue'

const matrixA = ref('')
const vectorB = ref('')
const initialGuess = ref('')
const epsilon = ref(0.0001)
const maxIter = ref(1000)
const error = ref('')
const result = ref('')

const handleSubmit = () => {
  error.value = ''
  
  // Валидация
  if (!matrixA.value.trim() || !vectorB.value.trim()) {
    error.value = 'заполните все поля'
    return
  }
  
  // значения по умолчанию если не заданы
  const eps = epsilon.value || 0.0001
  const maxIterations = maxIter.value || 1000
  
  // Валидация параметров
  if (eps <= 0) {
    error.value = 'eps должен быть больше 0'
    return
  }
  
  if (maxIterations <= 0 || !Number.isInteger(Number(maxIterations))) {
    error.value = 'предел итераций должен быть целым числом!!'
    return
  }
  
  const A = matrixA.value
    .trim()
    .split('\n')
    .map(row => row.trim().split(/\s+/).map(Number))
    .filter(row => row.length > 0)
  
  const b = vectorB.value
    .trim()
    .split('\n')
    .map(val => Number(val.trim()))
    .filter(val => !isNaN(val))
  
  if (A.length === 0 || b.length === 0) {
    error.value = 'некорректный формат данных'
    return
  }
  
  // все строки матрицы имеют одинаковую длину
  const firstRowLength = A[0].length
  const allRowsSameLength = A.every(row => row.length === firstRowLength)
  
  if (!allRowsSameLength) {
    error.value = 'все строки матрицы должны иметь одинаковое количество элементов'
    return
  }
  
  // матрица квадратная
  const n = A.length
  const m = A[0].length
  
  if (n !== m) {
    error.value = `матрица должна быть КВАДРАТНОЙ (сейчас ${n}x${m})`
    return
  }
  
  // размер вектора совпадает с размером матрицы
  if (b.length !== n) {
    error.value = `размер вектора b (${b.length}) должен совпадать с размером матрицы (${n}x${n})`
    return
  }
  
  // initial guess если задан
  let initialGuessArray = null
  if (initialGuess.value.trim()) {
    initialGuessArray = initialGuess.value
      .trim()
      .split('\n')
      .map(val => Number(val.trim()))
      .filter(val => !isNaN(val))
    
    if (initialGuessArray.length !== n) {
      error.value = `размер начального приближения (${initialGuessArray.length}) должен совпадать с размером матрицы (${n}x${n})`
      return
    }
  }
  
  result.value = 'Секунду...'

  const payload = {
    matrix: A,
    vector: b,
    epsilon: eps,
    max_iter: maxIterations
  }
  
  if (initialGuessArray) {
    payload.initial_guess = initialGuessArray
  }

  api.post('/api/solve', payload)
    .then(response => {
      // result.value = JSON.stringify(response.data, null, 2)
      processResult(response.data)
    })
    .catch(err => {
      error.value = (err.response?.data?.error || 'Ошибка запроса') + ": " + err.response?.data?.details
      result.value = ''
    })
}

async function processResult(resJson) {
  let x = resJson.solution
  let errors = resJson.errors
  let iters = resJson.iterations
  let norm = resJson.norm_of_matrix
  let messages = resJson.messages

  console.log('Response from server:', resJson)
  console.log('Norm of C:', norm)

  let res_str = '> найденное решение:\n'
  for (let i = 0; i < x.length; i++) {
    res_str += `x_${i + 1} = ${x[i]}\n`
  }
  res_str += `> количество итераций: ${iters}\n`
  res_str += `> погрешности на каждой итерации:\n`
  for (let i = 0; i < errors.length; i++) {
    res_str += `${i + 1}: ${errors[i]}\n`
  }
  res_str += `> норма приведенной матрицы C: ${norm}\n`

  let val = ''
  messages.forEach(msg => {
    val += msg
    val += '; '
  });
  error.value = val

  result.value = res_str
}

const handleFileUpload = (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const content = e.target.result
      const lines = content.trim().split('\n').filter(line => line.trim())
      
      const matrixRows = []
      const vectorValues = []
      
      for (const line of lines) {
        const values = line.trim().split(/\s+/).map(Number)
        if (values.length < 2) continue
        
        const bValue = values[values.length - 1]
        vectorValues.push(bValue)
        
        const matrixRow = values.slice(0, -1)
        matrixRows.push(matrixRow.join(' '))
      }
      
      matrixA.value = matrixRows.join('\n')
      vectorB.value = vectorValues.join('\n')
      error.value = ''
      
    } catch (err) {
      error.value = 'ридинг эрор: ' + err.message
    }
  }
  
  reader.onerror = () => {
    error.value = 'не удалось прочитать файл'
  }
  
  reader.readAsText(file)
  event.target.value = ''
}

const clearError = () => {
  error.value = ''
}
</script>

<template>
  <div class="lab-view">
    <h2 class="lab-heading">Решение СЛАУ методом Гаусса-Зейделя</h2>
    
    <Alert 
      v-if="error" 
      :message="error" 
      type="error" 
      @close="clearError" 
    />
    
    <div class="params-container">
      <div class="param-group">
        <label for="epsilon" class="form-label">Точность (eps):</label>
        <input 
          id="epsilon"
          v-model.number="epsilon" 
          type="number" 
          step="0.0001"
          class="form-input form-input-small"
          placeholder="0.0001"
        />
      </div>
      
      <div class="param-group">
        <label for="maxIter" class="form-label">Макс. итераций:</label>
        <input 
          id="maxIter"
          v-model.number="maxIter" 
          type="number" 
          step="1"
          class="form-input form-input-small"
          placeholder="1000"
        />
      </div>
    </div>
    
    <div class="form-group">
      <label class="form-label">Введите систему Ax = b:</label>
      <div class="matrix-input-container">
        <div class="matrix-input-col matrix-a">
          <label for="matrixA" class="input-sublabel">Матрица A</label>
          <textarea 
            id="matrixA"
            v-model="matrixA" 
            class="form-textarea"
            rows="8"
            placeholder="Коэффициенты матрицы A..."
          ></textarea>
        </div>
        <div class="matrix-input-col vector-b">
          <label for="vectorB" class="input-sublabel">Вектор b</label>
          <textarea 
            id="vectorB"
            v-model="vectorB" 
            class="form-textarea"
            rows="8"
            placeholder="Вектор b..."
          ></textarea>
        </div>
        <div class="matrix-input-col initial-guess">
          <label for="initialGuess" class="input-sublabel">Начальное приближение</label>
          <textarea 
            id="initialGuess"
            v-model="initialGuess" 
            class="form-textarea"
            rows="8"
            placeholder="Опционально..."
          ></textarea>
        </div>
      </div>
    </div>
    
    <div class="button-group">
      <button @click="handleSubmit" class="btn btn-primary">
        Отправить
      </button>
      
      <label class="btn btn-secondary" for="fileInput">
        Приложить файл
      </label>
      <input 
        id="fileInput"
        type="file" 
        accept=".txt"
        @change="handleFileUpload"
        style="display: none;"
      />
    </div>
    
    <div v-if="result" class="output">
      <div class="output-label">Вывод:</div>
      <pre class="output-value">{{ result }}</pre>
    </div>
  </div>
</template>

<style scoped lang="scss">
@import '@/styles/variables.scss';

.lab-view {
  padding: $spacing-md 0;
}

.lab-heading {
  font-family: $mono-font;
  font-size: $font-size-heading;
  color: $text-primary;
  margin: 0 0 $spacing-xl 0;
  font-weight: normal;
  border-bottom: $border-width solid $border-color;
  padding-bottom: $spacing-md;
}

.matrix-input-container {
  display: flex;
  gap: $spacing-md;
}

.matrix-input-col {
  display: flex;
  flex-direction: column;
}

.matrix-a {
  flex: 2;
}

.vector-b {
  flex: 1;
}

.initial-guess {
  flex: 1;
}

.input-sublabel {
  font-family: $mono-font;
  font-size: $font-size-small;
  color: $text-secondary;
  margin-bottom: $spacing-xs;
}

.params-container {
  display: flex;
  gap: $spacing-xl;
  margin-bottom: $spacing-lg;
}

.param-group {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  
  .form-label {
    margin-bottom: 0;
    white-space: nowrap;
  }
}

.form-input-small {
  max-width: 120px;
}

.button-group {
  display: flex;
  gap: $spacing-md;
  align-items: center;
}

.output {
  margin-top: $spacing-xl;
}
</style>
