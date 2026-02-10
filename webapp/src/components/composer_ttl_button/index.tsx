import React, {useState, useEffect} from 'react';
import {Store} from 'redux';
import {GlobalState} from '@mattermost/client';
import TTLMenu from './ttl_menu';
import './styles.css';

interface ComposerTTLButtonProps {
    store: Store<GlobalState>;
}

interface DurationOption {
    label: string;
    value: string;
}

const durations: DurationOption[] = [
    {label: '5 minutes', value: '5m'},
    {label: '15 minutes', value: '15m'},
    {label: '1 hour', value: '1h'},
    {label: '1 day', value: '1d'},
];

const ComposerTTLButton: React.FC<ComposerTTLButtonProps> = (props) => {
    const {store} = props;
    const [showMenu, setShowMenu] = useState(false);
    const [selectedDuration, setSelectedDuration] = useState<string | null>(null);

    useEffect(() => {
        const state = store.getState();
        const pluginState = (state as any).plugins?.['com.fambear.expiring-messages'];
        if (pluginState?.selectedDuration) {
            setSelectedDuration(pluginState.selectedDuration);
        }
    }, [store]);

    const handleSelect = (duration: string) => {
        setSelectedDuration(duration);
        setShowMenu(false);

        (window as any).setSelectedTTLDuration = duration;
    };

    const handleClear = () => {
        setSelectedDuration(null);
        setShowMenu(false);
        delete (window as any).setSelectedTTLDuration;
    };

    const getSelectedLabel = () => {
        if (!selectedDuration) return '';
        const found = durations.find(d => d.value === selectedDuration);
        return found ? found.label : selectedDuration;
    };

    return (
        <div className="composer-ttl-container">
            <button
                className={`composer-ttl-button ${selectedDuration ? 'active' : ''}`}
                onClick={() => setShowMenu(!showMenu)}
                title={selectedDuration ? `TTL: ${getSelectedLabel()}` : 'Set message expiration'}
            >
                <span className="flame-icon">🔥</span>
            </button>
            {showMenu && (
                <>
                    <div className="ttl-menu-overlay" onClick={() => setShowMenu(false)} />
                    <TTLMenu
                        durations={durations}
                        selected={selectedDuration}
                        onSelect={handleSelect}
                        onClear={handleClear}
                    />
                </>
            )}
        </div>
    );
};

export default ComposerTTLButton;
