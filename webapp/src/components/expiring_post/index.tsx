import React, {useState, useEffect} from 'react';
import './styles.css';

interface ExpiringPostProps {
    post: {
        id: string;
        message: string;
        user_id: string;
        create_at: number;
        props?: {
            ttl?: {
                enabled: boolean;
                duration: string;
                expires_at: number;
            };
        };
    };
    theme: any;
}

const formatTimeRemaining = (ms: number): string => {
    if (ms <= 0) return 'Expiring...';
    
    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    
    if (days > 0) {
        return `${days}d ${hours % 24}h`;
    }
    if (hours > 0) {
        return `${hours}h ${minutes % 60}m`;
    }
    if (minutes > 0) {
        return `${minutes}m ${seconds % 60}s`;
    }
    return `${seconds}s`;
};

const ExpiringPost: React.FC<ExpiringPostProps> = ({post, theme}) => {
    const [timeRemaining, setTimeRemaining] = useState<number>(0);
    
    const expiresAt = post.props?.ttl?.expires_at || 0;
    
    useEffect(() => {
        const updateTime = () => {
            const remaining = expiresAt - Date.now();
            setTimeRemaining(remaining);
        };
        
        updateTime();
        const interval = setInterval(updateTime, 1000);
        
        return () => clearInterval(interval);
    }, [expiresAt]);
    
    const isExpired = timeRemaining <= 0;
    const isUrgent = timeRemaining > 0 && timeRemaining < 60000; // Less than 1 minute
    
    return (
        <div className={`expiring-post ${isExpired ? 'expired' : ''} ${isUrgent ? 'urgent' : ''}`}>
            <div className="expiring-post-header">
                <span className="expiring-post-countdown">
                    🔥 {formatTimeRemaining(timeRemaining)}
                </span>
            </div>
            <div className="expiring-post-message">
                {post.message}
            </div>
        </div>
    );
};

export default ExpiringPost;
