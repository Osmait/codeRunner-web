import { useEffect, useState } from "react";
import "./App.css";
import { Editor } from "@monaco-editor/react";

function App() {
  const [code, setCode] = useState<string>("");
  const [output, setOutput] = useState("");
  const [socket, setSocket] = useState<WebSocket>();

  const [language, setLanguage] = useState("javascript"); // Lenguaje por defecto

  const handleEditorChange = (code: string | undefined) => {
    if (!code) return;
    setCode(code);
  };
  const handleLanguageChange = (event) => {
    setLanguage(event.target.value);
  };

  useEffect(() => {
    const ws = new WebSocket("ws://127.0.0.1:8080/ws");

    ws.onopen = () => {
      console.log("Conectado al servidor WebSocket");
    };

    ws.onmessage = (event) => {
      console.log("Mensaje recibido:", event);
      const response = JSON.parse(event.data);
      setOutput(response.result);
    };

    ws.onerror = (error) => {
      console.error("Error en WebSocket:", error);
    };

    ws.onclose = () => {
      console.log("Conexión WebSocket cerrada");
    };

    setSocket(ws);

    return () => ws.close();
  }, []);
  const executeCommand = () => {
    if (socket) {
      const message = JSON.stringify({
        lang: language, // Puedes cambiar esto si es necesario
        code: code,
      });

      console.log("Enviando comando:", message);
      socket.send(message);
    }
  };

  console.log(language);
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        height: "90vh",
        gap: "10px",
      }}
    >
      {/* Contenedor para el editor */}
      <div style={{ display: "flex", justifyContent: "space-between" }}>
        <select
          value={language}
          onChange={handleLanguageChange}
          style={{ padding: "8px", borderRadius: "5px", width: "200px" }}
        >
          <option value="javascript">JavaScript</option>
          <option value="python">Python</option>
        </select>

        <div>
          <button onClick={() => executeCommand()}>Ejecutar Código</button>
        </div>
      </div>

      <div
        style={{
          flexGrow: 1,
          flexBasis: "60%",
          display: "flex",
          flexDirection: "column",
          borderRadius: "10px", // Borde redondeado
          overflow: "hidden", // Asegura que el editor no sobresalga del contenedor
        }}
      >
        <Editor
          height="100%"
          width="100%" // El ancho será 100% del contenedor padre
          defaultLanguage="javascript"
          theme="vs-dark"
          defaultValue={code}
          onChange={handleEditorChange}
          options={{
            minimap: { enabled: false }, // Desactiva el minim
            fontSize: 16,
          }}
        />
      </div>

      {/* Contenedor para la terminal */}
      <div
        style={{
          flexGrow: 1,
          flexBasis: "40%",
          display: "flex",
          flexDirection: "column",
          gap: "10px",
        }}
      >
        <div
          style={{
            backgroundColor: "#1e1e1e", // Fondo oscuro
            color: "#00ff00", // Texto verde estilo terminal
            borderRadius: "10px",
            padding: "10px",
            height: "100%",
            overflow: "auto",
            fontFamily: "monospace", // Fuente monoespaciada
            whiteSpace: "pre-wrap", // Mantiene saltos de línea
            flexGrow: 1, // Hace que el contenedor de la terminal ocupe todo el espacio disponible
          }}
        >
          {output}
        </div>
      </div>
    </div>
  );
}
export default App;
