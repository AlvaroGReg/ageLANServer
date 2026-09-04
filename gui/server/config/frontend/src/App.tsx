import { useEffect, useRef, useState } from 'react';
import { FluentProvider, webDarkTheme } from '@fluentui/react-components';
import { ConfigFile, ConfigService, Configuration, ValidationError } from '../bindings/changeme';
import { AppConfig } from './types';
import { Screen1_FileLoad } from './components/Screen1_FileLoad';
import { Screen2_ConfigEditor } from './components/Screen2_ConfigEditor';

function toEditorConfig(config: Configuration): AppConfig {
    const game = (hosts: string[] | null) => ({ Hosts: hosts ?? [] });
    return {
        Log: config.Log,
        GeneratePlatformUserId: config.GeneratePlatformUserId,
        Authentication: config.Authentication as AppConfig['Authentication'],
        Games: {
            Enabled: config.Games.Enabled ?? [],
            age1: game(config.Games.age1.Hosts),
            age2: game(config.Games.age2.Hosts),
            age3: game(config.Games.age3.Hosts),
            age4: game(config.Games.age4.Hosts),
            athens: game(config.Games.athens.Hosts),
        },
        Announcement: config.Announcement,
    };
}

export function App() {
    const [loadedFile, setLoadedFile] = useState<ConfigFile | null>(null);
    const [errorMsg, setErrorMsg] = useState<string | null>(null);
    const [validationErrors, setValidationErrors] = useState<ValidationError[]>([]);
    const hasStartedLoading = useRef(false);

    const openConfigFile = () => {
        setErrorMsg(null);
        ConfigService.OpenConfigFile()
            .then((file) => {
                if (!file) return;
                if (file.validationErrors?.length) {
                    setValidationErrors(file.validationErrors);
                    setLoadedFile(null);
                    return;
                }
                setValidationErrors([]);
                setLoadedFile(file);
            })
            .catch((error: Error) => setErrorMsg(error.message));
    };

    const selectConfigFile = () => {
        setErrorMsg(null);
        ConfigService.SelectConfigFile()
            .then((file) => {
                if (!file) return;
                if (file.validationErrors?.length) {
                    setValidationErrors(file.validationErrors);
                    setLoadedFile(null);
                    return;
                }
                setValidationErrors([]);
                setLoadedFile(file);
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
                    {loadedFile ? (
                        <Screen2_ConfigEditor
                            initialConfig={toEditorConfig(loadedFile.config!)}
                            originalContent={loadedFile.content}
                            fileName={loadedFile.path}
                            onBack={() => setLoadedFile(null)}
                        />
                    ) : (
                        <Screen1_FileLoad
                            onSelectConfig={selectConfigFile}
                            errorMsg={errorMsg || undefined}
                            validationErrors={validationErrors}
                        />
                    )}
                </main>

                {/* WinUI 3 Footer bar */}
                <footer className="win-footer">
                    <div className="win-footer-brand">
                        <span className="win-footer-dot" />
                        <span>ageLANServer Configurator • Fluent UI React v9</span>
                    </div>
                    <div className="win-footer-info">
                        <span>{loadedFile ? 'Phase 2 • Configuration loaded' : 'Phase 2 • Load configuration'}</span>
                    </div>
                </footer>
            </div>
        </FluentProvider>
    );
}

export default App;
