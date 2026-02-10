import React, {useState, useEffect} from 'react';
import {formatTimeRemaining} from '../../utils/timer';
import './styles.css';

interface PostProps {
    post: {
        props?: {
            ttl?: {
                enabled: boolean;
                expires_at: number;
                duration: string;
            };
        };
    };
}

const PostTTLIndicator: React.FC<PostProps> = (props) => {
    const {post} = props;
    const ttl = post.props?.ttl;

    const [timeRemaining, setTimeRemaining] = useState<string | null>(null);
    const [isExpired, setIsExpired] = useState(false);

    useEffect(() => {
        if (!ttl?.enabled || !ttl?.expires_at) {
            return;
        }

        const updateTimer = () => {
            const now = Date.now();
            const remaining = ttl.expires_at - now;

            if (remaining <= 0) {
                setTimeRemaining('00:00');
                setIsExpired(true);
            } else {
                setTimeRemaining(formatTimeRemaining(remaining));
                setIsExpired(false);
            }
        };

        updateTimer();
        const interval = setInterval(updateTimer, 1000);

        return () => clearInterval(interval);
    }, [ttl]);

    if (!ttl?.enabled) {
        return null;
    }

    return (
        <span className={`ttl-indicator ${isExpired ? 'expired' : ''}`}>
            <span className="ttl-icon">🔥</span>
            {timeRemaining && (
                <span className="ttl-countdown">{timeRemaining}</span>
            )}
        </span>
    );
};

export default PostTTLIndicator;
