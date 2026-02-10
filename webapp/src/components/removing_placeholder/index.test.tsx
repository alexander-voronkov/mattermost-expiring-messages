import { describe, it, expect, vi } from 'vitest';
import { render, screen, cleanup } from '@testing-library/react';
import RemovingPlaceholder from './index';

describe('RemovingPlaceholder', () => {
  afterEach(() => {
    cleanup();
    vi.useRealTimers();
  });

  describe('when TTL is not enabled', () => {
    it('should render nothing when post has no props', () => {
      const { container } = render(<RemovingPlaceholder post={{}} />);
      expect(container.firstChild).toBeNull();
    });

    it('should render nothing when post has props but no TTL', () => {
      const { container } = render(<RemovingPlaceholder post={{ props: {} }} />);
      expect(container.firstChild).toBeNull();
    });

    it('should render nothing when TTL exists but is disabled', () => {
      const post = {
        props: {
          ttl: {
            enabled: false,
            expires_at: Date.now() - 1000,
          },
        },
      };
      const { container } = render(<RemovingPlaceholder post={post} />);
      expect(container.firstChild).toBeNull();
    });

    it('should render nothing when TTL enabled is undefined', () => {
      const post = {
        props: {
          ttl: {
            expires_at: Date.now() - 1000,
          },
        },
      };
      const { container } = render(<RemovingPlaceholder post={post} />);
      expect(container.firstChild).toBeNull();
    });
  });

  describe('when TTL is enabled but not expired', () => {
    it('should render nothing for future expiration', () => {
      vi.useFakeTimers().setSystemTime(new Date('2024-01-01T00:00:00Z'));

      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: new Date('2024-01-01T00:05:00Z').getTime(),
          },
        },
      };

      const { container } = render(<RemovingPlaceholder post={post} />);
      expect(container.firstChild).toBeNull();
    });

    it('should render nothing for expiration far in future', () => {
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: Date.now() + 24 * 60 * 60 * 1000, // 1 day from now
          },
        },
      };

      const { container } = render(<RemovingPlaceholder post={post} />);
      expect(container.firstChild).toBeNull();
    });

    it('should render nothing 1ms before expiration', () => {
      vi.useFakeTimers().setSystemTime(new Date('2024-01-01T00:00:00Z'));

      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: new Date('2024-01-01T00:00:00.001Z').getTime(),
          },
        },
      };

      const { container } = render(<RemovingPlaceholder post={post} />);
      expect(container.firstChild).toBeNull();
    });
  });

  describe('when TTL is enabled and expired', () => {
    it('should render removing text for expired post', () => {
      vi.useFakeTimers().setSystemTime(new Date('2024-01-01T00:05:00Z'));

      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: new Date('2024-01-01T00:04:59Z').getTime(),
          },
        },
      };

      render(<RemovingPlaceholder post={post} />);

      expect(screen.getByText('removing...')).toBeInTheDocument();
    });

    it('should render removing text when expires_at is exactly now', () => {
      const now = Date.now();
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: now,
          },
        },
      };

      render(<RemovingPlaceholder post={post} />);

      expect(screen.getByText('removing...')).toBeInTheDocument();
    });

    it('should render removing text for post expired long ago', () => {
      vi.useFakeTimers().setSystemTime(new Date('2024-01-01T01:00:00Z'));

      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: new Date('2024-01-01T00:00:00Z').getTime(),
          },
        },
      };

      render(<RemovingPlaceholder post={post} />);

      expect(screen.getByText('removing...')).toBeInTheDocument();
    });
  });

  describe('component structure', () => {
    it('should render removing-placeholder class', () => {
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: Date.now() - 1000,
          },
        },
      };

      const { container } = render(<RemovingPlaceholder post={post} />);

      const placeholder = container.querySelector('.removing-placeholder');
      expect(placeholder).toBeInTheDocument();
    });

    it('should render removing-text class', () => {
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: Date.now() - 1000,
          },
        },
      };

      const { container } = render(<RemovingPlaceholder post={post} />);

      const text = container.querySelector('.removing-text');
      expect(text).toBeInTheDocument();
      expect(text).toHaveTextContent('removing...');
    });
  });

  describe('edge cases', () => {
    it('should handle zero expires_at', () => {
      // Unix epoch is definitely in the past
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: 0,
          },
        },
      };

      render(<RemovingPlaceholder post={post} />);

      expect(screen.getByText('removing...')).toBeInTheDocument();
    });

    it('should handle negative expires_at', () => {
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: -1000,
          },
        },
      };

      render(<RemovingPlaceholder post={post} />);

      expect(screen.getByText('removing...')).toBeInTheDocument();
    });

    it('should handle missing expires_at', () => {
      const post = {
        props: {
          ttl: {
            enabled: true,
          },
        },
      };

      const { container } = render(<RemovingPlaceholder post={post} />);

      // NaN comparison will be false, so should not render
      expect(container.firstChild).toBeNull();
    });
  });
});
