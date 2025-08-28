import Unauthorized from "../assets/Unauthorized.png";

function UnauthorizedPage() {
  return (
    <div className="flex flex-col items-center place-content-center font-stretch-50%">
      <h1 className="text-6xl">Počkej!</h1>
      <h2 className="text-4xl">Tady nemáš, co dělat! My si na tebe posvítíme!</h2>
      <div className="flex items-center gap-3">
        <img src={Unauthorized} alt="Unauthorized" className="w-90 h-60"></img>
      </div>
      <h3 className="text-2xl font-semibold">NEAUTORIZOVÁN</h3>
    </div>
  );
}

export default UnauthorizedPage;
