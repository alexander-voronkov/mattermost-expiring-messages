import PostTTLIndicator from './components/post_ttl_indicator';
import ComposerTTLButton from './components/composer_ttl_button';

declare global {
    interface Window {
        registerPlugin: (pluginId: string, plugin: any) => void;
        setSelectedTTLDuration: (duration: string | null) => void;
    }
}

let selectedTTLDuration: string | null = null;

window.setSelectedTTLDuration = (duration: string | null) => {
    selectedTTLDuration = duration;
};

export default class Plugin {
    initialize(registry: any) {
        registry.registerPostEditorActionComponent(ComposerTTLButton);
        
        // Add dropdown menu item to show TTL info
        if (registry.registerPostDropdownMenuAction) {
            registry.registerPostDropdownMenuAction(
                '🔥 View TTL',
                (postId: string) => {
                    // This will be shown when clicking the menu item
                    const post = (window as any).store?.getState()?.entities?.posts?.posts?.[postId];
                    const ttl = post?.props?.ttl;
                    if (ttl?.enabled && ttl?.expires_at) {
                        const remaining = ttl.expires_at - Date.now();
                        if (remaining > 0) {
                            const mins = Math.floor(remaining / 60000);
                            const secs = Math.floor((remaining % 60000) / 1000);
                            alert(`Message expires in ${mins}m ${secs}s`);
                        } else {
                            alert('Message has expired and will be deleted soon');
                        }
                    } else {
                        alert('This message does not have TTL enabled');
                    }
                },
                // Filter - only show for posts with TTL
                (postId: string) => {
                    const post = (window as any).store?.getState()?.entities?.posts?.posts?.[postId];
                    return !!post?.props?.ttl?.enabled;
                }
            );
        }

        registry.registerMessageWillBePostedHook((post: any) => {
            if (selectedTTLDuration) {
                const duration = selectedTTLDuration;
                // Reset immediately after capturing the value
                selectedTTLDuration = null;
                // Notify the button to reset its visual state
                window.dispatchEvent(new CustomEvent('ttl-reset'));
                
                const newProps = {
                    ...(post.props || {}),
                    ttl: {
                        enabled: true,
                        duration,
                    },
                };
                return {post: {...post, props: newProps}};
            }

            if (post.props?.ttl?.enabled === false) {
                const newProps = {...post.props};
                delete newProps.ttl;
                return {post: {...post, props: newProps}};
            }

            return {post};
        });
    }
}
