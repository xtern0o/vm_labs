<script setup>
import { ref } from 'vue'
import api from '@/services/api.js'
import Alert from '@/components/Alert.vue'

const epsilon = ref(0.001)
const n = ref(4)
const a = ref(0)
const b = ref(1)
const result = ref('')
const error = ref('')

const defFuncOptions = [
  { id: 1, label: 'Функция 1', tex: '-2x^3-3x^2+x+5', type: 'def' },
  { id: 2, label: 'Функция 2', tex: 'e^{0,1x}', type: 'def' },
  { id: 3, label: 'Функция 3', tex: '\\sin(\\cos(x)) + 2', type: 'def' },
  { id: 4, label: 'Функция 4', tex: '2 \\cdot \\ln x', type: 'def' },
]

const improperFuncOptions = [
  { id: 5, label: 'Несобеств. 1', tex: '\\frac{1}{\\sqrt{x}}', type: 'imp' },
  { id: 6, label: 'Несобеств. 2', tex: '\\frac{1}{1-x}', type: 'imp' },
  { id: 7, label: 'Несобеств. 3', tex: '\\frac{1}{\\sqrt{2x-x^{2}}}', type: 'imp' },
]

const funcOptions = [
  ...defFuncOptions.map(f => ({ ...f, globalType: 'def' })),
  ...improperFuncOptions.map(f => ({ ...f, globalType: 'imp' })),
]

const integralMethodOptions = [
  { id: 1, label: 'Метод левых прямоугольников' },
  { id: 2, label: 'Метод правых прямоугольников' },
  { id: 3, label: 'Метод центральных прямоугольников' },
  { id: 4, label: 'Метод трапеций' },
  { id: 5, label: 'Метод Симпсона' },
]


const selectedFunc = ref(null)
const selectedFuncType = ref(null)
const selectedMethod = ref(null)

function solveIntegral() {
  const eps = epsilon.value || 0.001
  const nOut = n.value || 4
  const method = selectedMethod.value

  // Найти выбранную функцию и её тип
  const funcObj = funcOptions.find(f => f.id === selectedFunc.value && f.globalType === selectedFuncType.value)
  if (!funcObj) {
    error.value = 'Не выбрана функция'
    return
  }

  const payload = {
    func_id: funcObj.id,
    a: a.value,
    b: b.value,
    eps: eps,
    n: nOut,
  }

  let url = ''
  if (funcObj.globalType === 'imp') {
    url = `/api/integral/improper/${method}`
    payload.func_id -= defFuncOptions.length
  } else {
    url = `/api/integral/${method}`
  }

  api.post(url, payload)
    .then(response => {
      processResult(response)
    })
    .catch(err => {
      error.value = (err.response?.data?.error || `Ошибка запроса (${err.response.status})`)
      result.value = ''
    })
}

async function processResult(response) {
    const data = response.data
    console.log(data)

    const ans = data.value
    const messages = data.messages
    const rungeR = data.runge_r
    const nRes = data.n

    let resStr = ''

    if (messages.length > 0) {
        resStr += `messages: [\n`
        messages.forEach(msg => {
            resStr += `  "${msg}",\n`
        });
        resStr += `]\n`
    }
    resStr += `> ответ: ${ans}\n`
    resStr += `> R: ${rungeR}\n`
    resStr += `> n, при котором достигнута необходимая точность: ${nRes}\n`

    result.value = resStr
}

</script>
<template>
  <div class="lab-view">
    <h2 class="lab-heading">Численное вычисление интегралов</h2>
    
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
        <label for="maxIter" class="form-label">n (начальное разбиение):</label>
        <input 
          id="n"
          v-model.number="n" 
          type="number" 
          step="1"
          class="form-input form-input-small"
          placeholder="4"
        />
      </div>
    </div>

    <div>
      <section class="choice-card">
        <div class="choice-grid">
            <div>
                <h3 class="choice-title">Доступные функции:</h3>
                <div class="radio-list">
                <label
                  v-for="option in funcOptions"
                  :key="option.globalType + '-' + option.id"
                  class="radio-item"
                >
                  <input
                    type="radio"
                    :value="option.id"
                    v-model="selectedFunc"
                    :name="'func-' + option.globalType"
                    @change="selectedFuncType = option.globalType"
                  />
                  <span class="radio-content">
                    <span class="radio-caption">{{ option.label }}</span>
                    <span class="math-inline" v-katex="option.tex"></span>
                    <!-- <span v-if="option.globalType === 'imp'" class="improper-label">(несобственный)</span> -->
                  </span>
                </label>
                </div>
            </div>

            <div>
                <h3 class="choice-title">Способ решения:</h3>
                <div class="radio-list">
                <label
                    v-for="method in integralMethodOptions"
                    :key="method.id"
                    class="radio-item"
                >
                    <input
                    v-model="selectedMethod"
                    :value="method.id"
                    type="radio"
                    name="integralMethod"
                    />
                    <span>{{ method.label }}</span>
                </label>
                </div>
            </div>
        </div>
            

        

        <h3 class="choice-subtitle">Интервал [a; b]</h3>
        <div class="inline-inputs">
          <div class="input-stack">
            <label for="intervalA" class="form-label">a:</label>
            <input
              id="intervalA"
              v-model.number="b"
              type="number"
              step="0.01"
              class="form-input form-input-small"
              placeholder="0"
            />
          </div>

          <div class="input-stack">
            <label for="intervalB" class="form-label">b:</label>
            <input
              id="intervalB"
              v-model.number="a"
              type="number"
              step="0.01"
              class="form-input form-input-small"
              placeholder="1"
            />
          </div>
        </div>

        
        

        <!-- <label class="btn btn-secondary" for="fileInputEq">
          Приложить файл
        </label>
        <input 
          id="fileInputEq"
          type="file" 
          accept=".csv"
          @change="handleFileUploadEquation"
          style="display: none;"
        /> -->
        <button
          class="btn btn-primary solve-btn"
          :disabled="selectedFunc === null || selectedMethod === null"
          @click="solveIntegral"
        >
          Решить
        </button>
      </section>

    </div>
    
    <div v-if="result" class="output">
      <div class="output-label">Вывод:</div>
      <pre class="output-value">{{ result }}</pre>
    </div>
  </div>
</template>