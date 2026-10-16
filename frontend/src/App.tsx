import { useCalculatorController } from './controllers/useCalculatorController.ts'
import { CalculatorView } from './views/CalculatorView.tsx'

function App() {
  const calculator = useCalculatorController()
  return <CalculatorView {...calculator} />
}

export default App
