<script setup>
import { ref } from 'vue'
import api from '@/services/api.js'
import Alert from '@/components/Alert.vue'
import VueECharts from 'vue-echarts'
import * as echarts from 'echarts/core'
import { LineChart, ScatterChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, TitleComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

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
    tex: '\\begin{cases}\\tan(xy+0.3)-x^2=0\\\\0.9x^2+2y^2-1=0\\end{cases}',
  },
  {
    id: 2,
    label: 'Система 2',
    tex: '\\begin{cases}x^2+y^2-1=0\\\\x-y=0\\end{cases}',
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

const solveEquation = () => {

  const eps = epsilon.value || 0.001
  const maxIterations = maxIter.value || 1000
	const method = selectedEquationMethod.value

  const payload = {
    equation_id: selectedEquation.value,
    a: intervalA.value,
    b: intervalB.value,
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
  console.log(resJson)
}

const clearError = () => {
  error.value = ''
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
