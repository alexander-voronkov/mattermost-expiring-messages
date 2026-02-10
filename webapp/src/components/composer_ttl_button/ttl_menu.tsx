import React from 'react';
import './styles.css';

interface DurationOption {
    label: string;
    value: string;
}

interface TTLMenuProps {
    durations: DurationOption[];
    selected: string | null;
    onSelect: (duration: string) => void;
    onClear: () => void;
}

const TTLMenu: React.FC<TTLMenuProps> = ({durations, selected, onSelect, onClear}) => {
    return (
        <div className="ttl-menu">
            <div className="ttl-menu-header">
                <span className="ttl-menu-title">Message expires in...</span>
            </div>
            <div className="ttl-menu-options">
                {durations.map((duration) => (
                    <div
                        key={duration.value}
                        className={`ttl-option ${selected === duration.value ? 'selected' : ''}`}
                        onClick={() => onSelect(duration.value)}
                    >
                        <span className="ttl-option-label">{duration.label}</span>
                        {selected === duration.value && (
                            <span className="ttl-option-check">✓</span>
                        )}
                    </div>
                ))}
                {selected && (
                    <div
                        className="ttl-option ttl-option-clear"
                        onClick={onClear}
                    >
                        <span className="ttl-option-label">Disable TTL</span>
                    </div>
                )}
            </div>
        </div>
    );
};

export default TTLMenu;
