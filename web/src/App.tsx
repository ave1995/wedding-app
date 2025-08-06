import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import HomePage from "./pages/HomePage";
import InvitePage from "./pages/InvitePage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/invite" element={<InvitePage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
