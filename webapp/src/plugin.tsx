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
        registry.registerPostActionComponent(PostTTLIndicator);
        registry.registerPostEditorActionComponent(ComposerTTLButton);

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
