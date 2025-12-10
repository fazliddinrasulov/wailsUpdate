import { useState } from "react";
import "./App.css";
import {
  SilentPrint,
  CheckUpdate,
  DownloadUpdate,
  InstallUpdate,
} from "../wailsjs/go/main/App";

function App() {
  const [status, setStatus] = useState("");
  const [updateInfo, setUpdateInfo] = useState<any>(null);

  const handlePrint = async () => {
    try {
      setStatus("Printing...");
      await SilentPrint("C:\\Users\\Public\\test.pdf");
      setStatus("Printed");
    } catch (err) {
      setStatus("Print failed");
    }
  };

  const checkForUpdate = async () => {
    setStatus("Checking for update...");
    const info = await CheckUpdate("1.0.0"); // your current version
    if (info) {
      setUpdateInfo(info);
      setStatus("Update available");
    } else {
      setStatus("No update");
      setUpdateInfo(null);
    }
  };

  const downloadUpdate = async () => {
    if (!updateInfo) return;
    setStatus("Downloading update...");

    const dest = "C:\\Users\\Public\\update.exe"; // temporary location

    try {
      await DownloadUpdate(updateInfo.url, dest);
      setStatus("Downloaded. Installing...");
      await InstallUpdate(dest);
      setStatus("Installer launched");
    } catch (err) {
      setStatus("Update failed");
    }
  };

  return (
    <div id="App" style={{ padding: 20, fontFamily: "sans-serif" }}>
      <h2>Wails Test Panel</h2>

      <div style={{ marginTop: 20 }}>
        <button className="btn" onClick={handlePrint}>
          Test Silent Print
        </button>
      </div>

      <div style={{ marginTop: 20 }}>
        <button className="btn" onClick={checkForUpdate}>
          Check GitHub Update
        </button>
      </div>

      {updateInfo && (
        <div style={{ marginTop: 20, padding: 10, border: "1px solid #ccc" }}>
          <p>New version: {updateInfo.version}</p>
          <p>Release date: {updateInfo.release_date}</p>
          <p>Change log: {updateInfo.changelog}</p>

          <button className="btn" onClick={downloadUpdate}>
            Download and Install
          </button>
        </div>
      )}

      <div style={{ marginTop: 30, fontSize: 14 }}>
        <strong>Status:</strong> {status}
      </div>
    </div>
  );
}

export default App;
