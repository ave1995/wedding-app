import "./App.css";
import Button, { ButtonTypeEnum } from "./components/Button";

function App() {
  const simulateAsync = () =>
    new Promise<void>((resolve) => setTimeout(resolve, 1000));

  return (
    <Button
      onClickAsync={simulateAsync}
      label="Ping"
      type={ButtonTypeEnum.Basic}
    ></Button>
  );
}

export default App;
