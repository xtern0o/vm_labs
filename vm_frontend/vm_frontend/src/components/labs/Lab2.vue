<script setup>
import { ref } from 'vue'
import api from '@/services/api.js'
import Alert from '@/components/Alert.vue'
import VueECharts from 'vue-echarts'
import * as echarts from 'echarts/core'
import { LineChart, ScatterChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, TitleComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import Papa from "papaparse"

echarts.use([LineChart, ScatterChart, GridComponent, TooltipComponent, TitleComponent, CanvasRenderer])

const selectedEquation = ref(null)
const selectedSystem = ref(null)
const selectedEquationMethod = ref(null)
const selectedSystemMethod = ref(1)
const epsilon = ref(0.0001)
const maxIter = ref(1000)
const intervalA = ref(null)
const intervalB = ref(null)
const systemX0 = ref(null)
const systemY0 = ref(null)
const error = ref('')
const result = ref('')

const equationOptions = [
  { id: 1, label: 'Уравнение 1', tex: 'x^2 + e^x - 10x\\sin x - 5x = 0' },
  { id: 2, label: 'Уравнение 2', tex: '\\ln x + (\\sin x)^3 = 0' },
  { id: 3, label: 'Уравнение 3', tex: '\\sqrt{x} - 0.07x^3 + 0.5x^2 - 5 = 0' },
  { id: 4, label: 'Уравнение 4', tex: 'x^3 + 4x^2 + 3x - 1 = 0' },
  { id: 5, label: 'Уравнение 5', tex: '-0.4x^3 + 2x^2 - 2 = 0' },
]

const systemOptions = [
  {
    id: 1,
    label: 'Система 1',
    tex: '\\begin{cases}cos y+x-1.5=0\\\\sin x+y-4=0\\end{cases}',
  },
  {
    id: 2,
    label: 'Система 2',
    tex: '\\begin{cases}x^2+y^2-1=0\\\\x^2-3y=0.5\\end{cases}',
  },
]

const equationMethodOptions = [
  { id: 1, label: 'Метод половинного деления' },
  { id: 2, label: 'Метод Ньютона' },
  { id: 3, label: 'Метод простой итерации' },
]

const systemMethodOptions = [
  { id: 1, label: 'Метод простых итераций' },
]

const chartOption = ref(null)
const showChart = ref(false)

function plotSelectedEquation(solution, steps) {
  if (selectedEquation.value == null || intervalA.value == null || intervalB.value == null) return
  const a = Number(intervalA.value)
  const b = Number(intervalB.value)
  if (isNaN(a) || isNaN(b) || a >= b) return
  let f
  switch (selectedEquation.value) {
    case 1:
      f = x => x ** 2 + Math.exp(x) - 10 * x * Math.sin(x) - 5 * x
      break
    case 2:
      f = x => Math.log(x) + Math.pow(Math.sin(x), 3)
      break
    case 3:
      f = x => Math.sqrt(x) - 0.07 * x ** 3 + 0.5 * x ** 2 - 5
      break
    case 4:
      f = x => x ** 3 + 4 * x ** 2 + 3 * x - 1
      break
    case 5:
      f = x => -0.4 * x ** 3 + 2 * x ** 2 - 2
      break
    default:
      f = x => 0
  }
  const N = 100
  const xs = Array.from({ length: N + 1 }, (_, i) => a + (b - a) * i / N)
  const ys = xs.map(f)

  // приближения (серые)
  const stepPoints = Array.isArray(steps)
    ? steps.map(pt => [pt.x, pt.y])
    : []
  // решение (красное)
  const solutionPoint = solution && typeof solution.x === 'number' && typeof solution.y === 'number'
    ? [[solution.x, solution.y]]
    : []

	console.log(stepPoints)
	console.log(solutionPoint)

  chartOption.value = {
    title: { text: 'График функции', left: 'center' },
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 40, bottom: 40 },
    xAxis: { type: 'value', name: 'x', min: a, max: b },
    yAxis: { type: 'value', name: 'f(x)' },
    series: [
      {
        type: 'line',
        data: xs.map((x, i) => [x, ys[i]]),
        smooth: true,
        showSymbol: false,
        name: 'f(x)'
      },
      stepPoints.length > 0 ? {
        type: 'scatter',
        data: stepPoints,
        symbolSize: 10,
        itemStyle: { color: '#888', opacity: 0.7 },
        name: 'steps',
        z: 3
      } : null,
      solutionPoint.length > 0 ? {
        type: 'scatter',
        data: solutionPoint,
        symbolSize: 14,
        itemStyle: { color: '#e74c3c' },
        name: 'solution',
        z: 4
      } : null
    ].filter(Boolean),
  }
  showChart.value = true
}

function plotSelectedSystem(solution, steps) {
  if (!selectedSystem.value) return

  let f1, f2, domain
  switch (selectedSystem.value) {
    case 1:
      f1 = (x, y) => Math.cos(y) + x - 1.5
      f2 = (x, y) => Math.sin(x) + y - 1

      domain = { xMin: -4, xMax: 4, yMin: -2, yMax: 6 }
      break
    case 2:
      f1 = (x, y) => x ** 2 + y ** 2 - 1
      f2 = (x, y) => x ** 2 - 3 * y - 0.5
      domain = { xMin: -2, xMax: 2, yMin: -2, yMax: 2 }
      break
    default:
      f1 = (x, y) => 0
      f2 = (x, y) => 0
      domain = { xMin: -2, xMax: 2, yMin: -2, yMax: 2 }
  }

  // строим level curves (контуры) этих функций = 0
  const xSamples = 200
  const ySamples = 200
  const xVals = Array.from({ length: xSamples + 1 }, (_, i) =>
    domain.xMin + (domain.xMax - domain.xMin) * i / xSamples)
  const yVals = Array.from({ length: ySamples + 1 }, (_, j) =>
    domain.yMin + (domain.yMax - domain.yMin) * j / ySamples)

  // строим lines с мелкими контурами
  // берём только те (x, y), где f меняет знак
  // sampling |f| < 0.05

  // найдем где f1(x, y)~0 и f2(x, y)~0
  function findZeroContours(f, threshold = 0.005) {
    const points = []
    for (let xi = 0; xi < xSamples; ++xi) {
      for (let yi = 0; yi < ySamples; ++yi) {
        const x = xVals[xi]
        const y = yVals[yi]
        const f0 = f(x, y)
        const f1_ = f(xVals[xi+1], y)
        const f2_ = f(x, yVals[yi+1])

        if (Math.abs(f0) < threshold) {
          points.push([x, y])
        } else if (f0 * f1_ < 0) {
          points.push([(x + xVals[xi+1]) / 2, y])
        } else if (f0 * f2_ < 0) {
          points.push([x, (y + yVals[yi+1]) / 2])
        }
      }
    }
    return points
  }

  const zeros_f1 = findZeroContours(f1)
  const zeros_f2 = findZeroContours(f2)

  // серенькие эстимейты
  const stepPoints = Array.isArray(steps)
    ? steps.map(pt => [pt.x, pt.y])
    : []
  // солюшин
  const solutionPoint =
    solution && typeof solution.x === 'number' && typeof solution.y === 'number'
      ? [[solution.x, solution.y]]
      : []

  chartOption.value = {
    title: { text: 'График системы', left: 'center' },
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 40, bottom: 40 },
    xAxis: { type: 'value', name: 'x', min: domain.xMin, max: domain.xMax },
    yAxis: { type: 'value', name: 'y', min: domain.yMin, max: domain.yMax },
    series: [
      // f1=0 (синий)
      {
        type: 'scatter',
        data: zeros_f1,
        symbolSize: 3,
        itemStyle: { color: '#005bea' }, // синий
        name: 'f₁(x, y) = 0',
        z: 2,
      },
      // f2=0 (зелёный)
      {
        type: 'scatter',
        data: zeros_f2,
        symbolSize: 3,
        itemStyle: { color: '#30b878' }, // зелёный
        name: 'f₂(x, y) = 0',
        z: 2,
      },
      // шаги (серые)
      stepPoints.length > 0
        ? {
            type: 'scatter',
            data: stepPoints,
            symbolSize: 10,
            itemStyle: { color: '#888', opacity: 0.7 },
            name: 'steps',
            z: 3,
          }
        : null,
      // решение (красное)
      solutionPoint.length > 0
        ? {
            type: 'scatter',
            data: solutionPoint,
            symbolSize: 14,
            itemStyle: { color: '#e74c3c' },
            name: 'solution',
            z: 4,
          }
        : null,
    ].filter(Boolean),
  }
  showChart.value = true
}

const solveEquation = () => {

  const eps = epsilon.value || 0.001
  const maxIterations = maxIter.value || 1000
	const method = selectedEquationMethod.value
  const intA = intervalA.value || 0
  const intB = intervalB.value || 1

  const payload = {
    equation_id: selectedEquation.value,
    a: intA,
    b: intB,
    eps: eps,
		max_iter: maxIterations
  }

  api.post(`/api/equation/${method}`, payload)
    .then(response => {
      processResultEquation(response.data)
    })
    .catch(err => {
      error.value = (err.response?.data?.error || 'Ошибка запроса')
      showChart.value = false
      result.value = ''
    })
}

const solveSystem = () => {
  const eps = epsilon.value || 0.001
  const maxIterations = maxIter.value || 1000
	const method = selectedSystemMethod.value
  const x0 = systemX0.value || 0
  const y0 = systemY0.value || 0

  const payload = {
    system_id: selectedSystem.value,
    x0: x0,
    y0: y0,
    eps: eps,
    max_iter: maxIterations,    
  }

  api.post(`/api/system/${method}`, payload)
    .then(response => {
      processResultSystem(response.data)
    })
    .catch(err => {
      error.value = (err.response?.data?.error || 'Ошибка запроса')
      showChart.value = false
      result.value = ''
    })
}

async function processResultEquation(resJson) {
	console.log(resJson)
	
	const ans = {x: resJson.result.solution, y: resJson.result.value}
	const messages = resJson.result.messages
	const arg_error = resJson.result.arg_error
	const steps = resJson.result.steps
  const iter = resJson.result.iterations

  let resStr = ""

  if (messages.length > 0) {
      resStr += `messages: [\n`
      messages.forEach(msg => {
        resStr += `  "${msg}",\n`
      });
      resStr += `]\n`
  }

  resStr += `> найденное решение: x = ${ans.x}\n`
  resStr += `> f(x) = ${ans.y}\n`
  resStr += `> количество итераций: ${iter}\n`
  
  resStr += `> шаги приближения: [\n`
  steps.forEach(step => {
    resStr += `  (${step.x}, ${step.y}),\n`
  })
  resStr += `]\n`
  resStr += `> погрешность аргумента |x_i - x_(i+1)| = ${arg_error}\n`


  result.value = resStr
	plotSelectedEquation(ans, steps)
}

async function processResultSystem(resJson) {
  const solution = resJson.solution
  const steps = resJson.steps
  const messages = resJson.messages
  const errors = resJson.errors
  const iters = resJson.iterations

  console.log(resJson)

  let resStr = ""

  if (messages.length > 0) {
      resStr += `messages: [\n`
      messages.forEach(msg => {
        resStr += `  "${msg}",\n`
      });
      resStr += `]\n`
  }
  resStr += `> найденное решение: (${solution.x}, ${solution.y})\n`
  resStr += `> количество итераций: ${iters}\n`
  resStr += `> шаги приближения: [\n`
  steps.forEach(step => {
    resStr += `  (${step.x}, ${step.y}),\n`
  })
  resStr += `]\n`
  resStr += `> вектор погрешностей: [\n`
  errors.forEach(err => {
    resStr += `  (${err.x}, ${err.y})\n`
  })
  resStr += `]\n`

  result.value = resStr

  plotSelectedSystem(solution, steps)
}

const clearError = () => {
  error.value = ''
}

// equationOptions, equationMethodOptions, systemOptions, systemMethodOptions
// предполагается, что объявлены во внешней области

function handleFileUploadEquation(event) {
  const file = event.target.files[0]
  if (!file) return

  Papa.parse(file, {
    header: true, // парсим как объекты, имена — из первой строки CSV
    skipEmptyLines: true,
    complete: (results) => {
      try {
        const data = results.data
        if (!data.length) throw new Error('Нет строк с данными')

        // Берем первую строку-объект
        const row = data[0]

        // Проверка наличия и валидности значений
        const fields = ['eps', 'maxIter', 'equationId', 'methodId', 'a', 'b']
        for (const f of fields) {
          if (!(f in row) || row[f] === "") {
            throw new Error(`В файле отсутствует обязательное поле "${f}"`)
          }
        }

        const [eps, maxIterVal, equationId, methodId, aVal, bVal] = [
          Number(row.eps), Number(row.maxIter), Number(row.equationId),
          Number(row.methodId), Number(row.a), Number(row.b)
        ]

        if (![eps, maxIterVal, equationId, methodId, aVal, bVal].every(v => !isNaN(v))) {
          throw new Error('Некорректный формат: одно из значений не является числом')
        }

        const foundEquation = equationOptions.find(opt => opt.id === equationId)
        if (!foundEquation) throw new Error(`Уравнение с id=${equationId} не найдено`)

        const foundMethod = equationMethodOptions.find(opt => opt.id === methodId)
        if (!foundMethod) throw new Error(`Метод с id=${methodId} не найден`)

        epsilon.value = eps
        maxIter.value = maxIterVal
        selectedEquation.value = equationId
        selectedEquationMethod.value = methodId
        intervalA.value = aVal
        intervalB.value = bVal
        error.value = ''
      } catch (err) {
        error.value = `Ошибка загрузки файла уравнения: ${err.message}`
      }
    },
    error: (err) => {
      error.value = `Ошибка чтения файла уравнения: ${err.message}`
    }
  })
}

function handleFileUploadSystem(event) {
  const file = event.target.files[0]
  if (!file) return

  Papa.parse(file, {
    header: true,
    skipEmptyLines: true,
    complete: (results) => {
      try {
        const data = results.data
        if (!data.length) throw new Error('Нет строк с данными')

        // Берем первую строку-объект
        const row = data[0]
        const fields = ['eps', 'maxIter', 'systemId', 'methodId', 'x0', 'y0']
        for (const f of fields) {
          if (!(f in row) || row[f] === "") {
            throw new Error(`В файле отсутствует обязательное поле "${f}"`)
          }
        }

        const [eps, maxIterVal, systemId, methodId, x0, y0] =
          [Number(row.eps), Number(row.maxIter), Number(row.systemId), Number(row.methodId), Number(row.x0), Number(row.y0)]
        if (![eps, maxIterVal, systemId, methodId, x0, y0].every(v => !isNaN(v))) {
          throw new Error('Некорректный формат: одно из значений не является числом')
        }

        const foundSystem = systemOptions.find(opt => opt.id === systemId)
        if (!foundSystem) throw new Error(`Система с id=${systemId} не найдена`)

        const foundMethod = systemMethodOptions.find(opt => opt.id === methodId)
        if (!foundMethod) throw new Error(`Метод с id=${methodId} не найден`)

        epsilon.value = eps
        maxIter.value = maxIterVal
        selectedSystem.value = systemId
        selectedSystemMethod.value = methodId
        systemX0.value = x0
        systemY0.value = y0
        error.value = ''
      } catch (err) {
        error.value = `Ошибка загрузки файла системы: ${err.message}`
      }
    },
    error: (err) => {
      error.value = `Ошибка чтения файла системы: ${err.message}`
    }
  })
}

</script>

<template>
  <div class="lab-view">
    <h2 class="lab-heading">Численное решение нелинейных уравнений и систем</h2>

    <Alert
      v-if="error"
      :message="error"
      type="error"
      @close="clearError"
    />

    <div class="params-container">
      <div class="param-group">
        <label for="epsilon" class="form-label">точность (eps):</label>
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
          <label for="maxIter" class="form-label">max iter:</label>
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

    <div class="choice-grid">
      <section class="choice-card">
        <h3 class="choice-title">Нелинейные уравнения</h3>

        <div class="radio-list">
          <label
            v-for="option in equationOptions"
            :key="option.id"
            class="radio-item"
          >
            <input
              v-model="selectedEquation"
              :value="option.id"
              type="radio"
              name="equation"
            />
            <span class="radio-content">
              <span class="radio-caption">{{ option.label }}</span>
              <span class="math-inline" v-katex="option.tex"></span>
            </span>
          </label>
        </div>

        <h3 class="choice-subtitle">Интервал [a; b]</h3>
        <div class="inline-inputs">
          <div class="input-stack">
            <label for="intervalA" class="form-label">a:</label>
            <input
              id="intervalA"
              v-model.number="intervalA"
              type="number"
              step="0.1"
              class="form-input form-input-small"
              placeholder="0"
            />
          </div>

          <div class="input-stack">
            <label for="intervalB" class="form-label">b:</label>
            <input
              id="intervalB"
              v-model.number="intervalB"
              type="number"
              step="0.1"
              class="form-input form-input-small"
              placeholder="1"
            />
          </div>
        </div>

        <h3 class="choice-subtitle">Способ решения</h3>
        <div class="radio-list">
          <label
            v-for="method in equationMethodOptions"
            :key="method.id"
            class="radio-item"
          >
            <input
              v-model="selectedEquationMethod"
              :value="method.id"
              type="radio"
              name="equationMethod"
            />
            <span>{{ method.label }}</span>
          </label>
        </div>

        <label class="btn btn-secondary" for="fileInputEq">
          Приложить файл
        </label>
        <input 
          id="fileInputEq"
          type="file" 
          accept=".csv"
          @change="handleFileUploadEquation"
          style="display: none;"
        />
        <button
          class="btn btn-primary solve-btn"
          :disabled="selectedEquation === null || selectedEquationMethod === null"
          @click="solveEquation"
        >
          Решить
        </button>
      </section>

      <section class="choice-card">
        <h3 class="choice-title">Системы нелинейных уравнений</h3>

        <div class="radio-list">
          <label
            v-for="option in systemOptions"
            :key="option.id"
            class="radio-item"
          >
            <input
              v-model="selectedSystem"
              :value="option.id"
              type="radio"
              name="system"
            />
            <span class="radio-content">
              <span class="radio-caption">{{ option.label }}</span>
              <span class="math-display" v-katex:display="option.tex"></span>
            </span>
          </label>
        </div>

        <h3 class="choice-subtitle">Начальное приближение</h3>
        <div class="inline-inputs">
          <div class="input-stack">
            <label for="systemX0" class="form-label">x_0:</label>
            <input
              id="systemX0"
              v-model.number="systemX0"
              type="number"
              step="0.1"
              class="form-input form-input-small"
              placeholder="0"
            />
          </div>

          <div class="input-stack">
            <label for="systemY0" class="form-label">y_0:</label>
            <input
              id="systemY0"
              v-model.number="systemY0"
              type="number"
              step="0.1"
              class="form-input form-input-small"
              placeholder="0"
            />
          </div>
        </div>

        <h3 class="choice-subtitle">Способ решения</h3>
        <div class="radio-list">
          <label
            v-for="method in systemMethodOptions"
            :key="method.id"
            class="radio-item"
          >
            <input
              v-model="selectedSystemMethod"
              :value="method.id"
              type="radio"
              name="systemMethod"
            />
            <span>{{ method.label }}</span>
          </label>
        </div>

        <label class="btn btn-secondary" for="fileInputSys">
          Приложить файл
        </label>
        <input 
          id="fileInputSys"
          type="file" 
          accept=".csv"
          @change="handleFileUploadSystem"
          style="display: none;"
        />
        <button
          class="btn btn-primary solve-btn"
          :disabled="selectedSystem === null || selectedSystemMethod === null"
          @click="solveSystem"
        >
          Решить
        </button>
      </section>
    </div>
    <div v-if="showChart" class="output">
      <div class="output-label">График функции:</div>
      <VueECharts :option="chartOption" style="height: 320px; width: 100%;" />
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

.choice-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: $spacing-lg;
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

.inline-inputs {
  display: flex;
  gap: $spacing-md;
  flex-wrap: wrap;
}

.input-stack {
  display: flex;
  align-items: center;
  gap: $spacing-sm;

  .form-label {
    margin-bottom: 0;
    white-space: nowrap;
  }
}

.choice-card {
  border: $border-width solid $border-color;
  background: $bg-secondary;
  padding: $spacing-lg;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.choice-title {
  margin: 0;
  font-size: $font-size-base;
  font-weight: normal;
  color: $text-primary;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.choice-subtitle {
  margin: $spacing-sm 0 0 0;
  font-size: $font-size-small;
  font-weight: normal;
  color: $text-secondary;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.radio-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

.radio-item {
  display: flex;
  align-items: flex-start;
  gap: $spacing-sm;
  color: $text-primary;
  cursor: pointer;
}

.radio-content {
  display: flex;
  flex-direction: column;
  gap: $spacing-xs;
}

.radio-caption {
  color: $text-primary;
}

.math-inline {
  color: $text-secondary;
}

.math-display {
  color: $text-secondary;
  padding: $spacing-xs 0;
}

:deep(.katex-display) {
  margin: 0;
}

.radio-item input[type='radio'] {
  appearance: none;
  width: 14px;
  height: 14px;
  border: $border-width solid $border-color;
  border-radius: 50%;
  background: $bg-primary;
  display: inline-block;
  position: relative;
}

.radio-item input[type='radio']:checked {
  border-color: $accent-color;
}

.radio-item input[type='radio']:checked::after {
  content: '';
  position: absolute;
  top: 3px;
  left: 3px;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: $accent-color;
}

.radio-item input[type='radio']:focus-visible {
  outline: 1px solid $accent-color;
  outline-offset: 2px;
}

.solve-btn {
  margin-top: auto;
}

.button-group {
  display: flex;
  gap: $spacing-md;
  align-items: center;
}

.output {
  margin-top: $spacing-xl;
  background: $bg-secondary;
  padding: $spacing-lg;
  border-radius: $border-radius;
  box-shadow: 0 2px 8px rgba(0,0,0,0.03);
}
.output-label {
  font-weight: bold;
  margin-bottom: $spacing-sm;
}

@media (max-width: 860px) {
  .choice-grid {
    grid-template-columns: 1fr;
  }
}
</style>
