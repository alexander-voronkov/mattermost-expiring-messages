import React from 'react';

interface RemovingPlaceholderProps {
    post: {
        props?: {
            ttl?: {
                enabled: boolean;
                expires_at: number;
            };
        };
    };
}

const RemovingPlaceholder: React.FC<RemovingPlaceholderProps> = (props) => {
    const {post} = props;
    const ttl = post.props?.ttl;

    if (!ttl?.enabled) {
        return null;
    }

    const isExpired = Date.now() >= ttl.expires_at;

    if (!isExpired) {
        return null;
    }

    return (
        <div className="removing-placeholder">
            <span className="removing-text">removing...</span>
        </div>
    );
};

export default RemovingPlaceholder;
