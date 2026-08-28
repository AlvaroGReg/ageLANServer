import { useEffect, useRef, useState } from 'react';
import { FluentProvider, webDarkTheme } from '@fluentui/react-components';
import { AppConfig } from './types';
import { parseTOML } from './tomlUtils';
import { ConfigService } from '../bindings/changeme';
import { Screen1_FileLoad } from './components/Screen1_FileLoad';
import { Screen2_ConfigEditor } from './components/Screen2_ConfigEditor';

export function App() {
  const [loadedConfig, setLoadedConfig] = useState<AppConfig | null>(null);
  const [fileName, setFileName] = useState<string>('');
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const hasStartedLoading = useRef(false);

  const openConfigFile = () => {
    setErrorMsg(null);
    ConfigService.OpenConfigFile()
      .then((file) => {
        if (!file) return;
        setLoadedConfig(parseTOML(file.content));
        setFileName(file.path);
      })
      .catch((error: Error) => setErrorMsg(error.message));
  };

  useEffect(() => {
    if (hasStartedLoading.current) return;
    hasStartedLoading.current = true;
    openConfigFile();
  }, []);

  // Custom transparent background theme override so Wails 3 backdrop effect is visible
  const customFluentTheme = {
    ...webDarkTheme,
    colorNeutralBackground1: 'transparent',
    colorNeutralBackground2: 'rgba(255, 255, 255, 0.035)',
    colorNeutralBackground3: 'rgba(255, 255, 255, 0.06)',
    colorBrandBackground: '#60cdff',
    colorBrandBackgroundHover: '#4cc2ff',
    colorBrandForeground1: '#60cdff',
  };

  return (
    <FluentProvider theme={customFluentTheme} style={{ background: 'transparent', minHeight: '100vh' }}>
      <div className="win-app-root">
        {/* Wails 3 Backdrop Layer */}
        <div className="win-bg-backdrop" />

        <main className="win-main-content">
          {loadedConfig ? (
            <Screen2_ConfigEditor
              initialConfig={loadedConfig}
              fileName={fileName}
              onBack={() => setLoadedConfig(null)}
            />
          ) : (
            <Screen1_FileLoad onOpenConfig={openConfigFile} errorMsg={errorMsg || undefined} onError={setErrorMsg} />
          )}
        </main>

        {/* WinUI 3 Footer bar */}
        <footer className="win-footer">
          <div className="win-footer-brand">
            <span className="win-footer-dot" />
            <span>ageLANServer Configurator • Fluent UI React v9</span>
          </div>
          <div className="win-footer-info">
            <span>Screen: {loadedConfig ? '2 / 2 (Editor)' : '1 / 2 (Load)'}</span>
          </div>
        </footer>
      </div>
    </FluentProvider>
  );
}

export default App;
