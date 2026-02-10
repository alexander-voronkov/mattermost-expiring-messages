import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, cleanup } from '@testing-library/react';
import PostTTLIndicator from './index';

describe('PostTTLIndicator', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    cleanup();
    vi.useRealTimers();
  });

  describe('when TTL is not enabled', () => {
    it('should render nothing when post has no props', () => {
      const { container } = render(<PostTTLIndicator post={{}} />);
      expect(container.firstChild).toBeNull();
    });

    it('should render nothing when post has props but no TTL', () => {
      const { container } = render(<PostTTLIndicator post={{ props: {} }} />);
      expect(container.firstChild).toBeNull();
    });

    it('should render nothing when TTL exists but is disabled', () => {
      const post = {
        props: {
          ttl: {
            enabled: false,
            expires_at: Date.now() + 300000,
            duration: '5m',
          },
        },
      };
      const { container } = render(<PostTTLIndicator post={post} />);
      expect(container.firstChild).toBeNull();
    });

    it('should render nothing when TTL enabled is undefined', () => {
      const post = {
        props: {
          ttl: {
            expires_at: Date.now() + 300000,
            duration: '5m',
          },
        },
      };
      const { container } = render(<PostTTLIndicator post={post} />);
      expect(container.firstChild).toBeNull();
    });
  });

  describe('when TTL is enabled', () => {
    it('should render the flame icon and countdown', () => {
      const expiresAt = Date.now() + 300000; // 5 minutes from now
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: expiresAt,
            duration: '5m',
          },
        },
      };

      render(<PostTTLIndicator post={post} />);

      expect(screen.getByText('🔥')).toBeInTheDocument();
      expect(screen.getByText('05:00')).toBeInTheDocument();
    });

    it('should render countdown for 15 minutes', () => {
      const expiresAt = Date.now() + 15 * 60 * 1000;
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: expiresAt,
            duration: '15m',
          },
        },
      };

      render(<PostTTLIndicator post={post} />);

      expect(screen.getByText('🔥')).toBeInTheDocument();
      expect(screen.getByText('15:00')).toBeInTheDocument();
    });

    it('should render countdown for 1 hour', () => {
      const expiresAt = Date.now() + 60 * 60 * 1000;
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: expiresAt,
            duration: '1h',
          },
        },
      };

      render(<PostTTLIndicator post={post} />);

      expect(screen.getByText('🔥')).toBeInTheDocument();
      expect(screen.getByText('1:00:00')).toBeInTheDocument();
    });

    it('should render countdown for 1 day', () => {
      const expiresAt = Date.now() + 24 * 60 * 60 * 1000;
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: expiresAt,
            duration: '1d',
          },
        },
      };

      render(<PostTTLIndicator post={post} />);

      expect(screen.getByText('🔥')).toBeInTheDocument();
      expect(screen.getByText('24:00:00')).toBeInTheDocument();
    });
  });

  describe('countdown timer behavior', () => {
    it('should update countdown every second', () => {
      const expiresAt = Date.now() + 65000; // 1 minute 5 seconds
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: expiresAt,
            duration: '5m',
          },
        },
      };

      render(<PostTTLIndicator post={post} />);

      expect(screen.getByText('01:05')).toBeInTheDocument();

      // Advance time by 1 second
      vi.advanceTimersByTime(1000);

      expect(screen.getByText('01:04')).toBeInTheDocument();
    });

    it('should show 00:00 when expired', () => {
      const expiresAt = Date.now() - 1000; // 1 second ago
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: expiresAt,
            duration: '5m',
          },
        },
      };

      render(<PostTTLIndicator post={post} />);

      expect(screen.getByText('00:00')).toBeInTheDocument();
    });

    it('should add expired class when time is up', () => {
      const expiresAt = Date.now() - 1000;
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: expiresAt,
            duration: '5m',
          },
        },
      };

      const { container } = render(<PostTTLIndicator post={post} />);

      const ttlIndicator = container.querySelector('.ttl-indicator');
      expect(ttlIndicator).toHaveClass('expired');
    });

    it('should not add expired class when time remains', () => {
      const expiresAt = Date.now() + 60000;
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: expiresAt,
            duration: '5m',
          },
        },
      };

      const { container } = render(<PostTTLIndicator post={post} />);

      const ttlIndicator = container.querySelector('.ttl-indicator');
      expect(ttlIndicator).not.toHaveClass('expired');
    });
  });

  describe('edge cases', () => {
    it('should handle missing expires_at', () => {
      const post = {
        props: {
          ttl: {
            enabled: true,
            duration: '5m',
          },
        },
      };

      const { container } = render(<PostTTLIndicator post={post} />);

      // Should render the flame icon but no countdown
      expect(screen.getByText('🔥')).toBeInTheDocument();
      expect(container.querySelector('.ttl-countdown')).toBeNull();
    });

    it('should handle zero expires_at', () => {
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: 0,
            duration: '5m',
          },
        },
      };

      render(<PostTTLIndicator post={post} />);

      expect(screen.getByText('🔥')).toBeInTheDocument();
      expect(screen.getByText('00:00')).toBeInTheDocument();
    });

    it('should handle negative expires_at', () => {
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: -1000,
            duration: '5m',
          },
        },
      };

      render(<PostTTLIndicator post={post} />);

      expect(screen.getByText('🔥')).toBeInTheDocument();
      expect(screen.getByText('00:00')).toBeInTheDocument();
    });

    it('should handle very large expires_at', () => {
      const farFuture = Date.now() + 365 * 24 * 60 * 60 * 1000; // 1 year
      const post = {
        props: {
          ttl: {
            enabled: true,
            expires_at: farFuture,
            duration: '1d',
          },
        },
      };

      render(<PostTTLIndicator post={post} />);

      expect(screen.getByText('🔥')).toBeInTheDocument();
      // Should show hours in the countdown
      expect(screen.queryByText(/\d+:\d{2}:\d{2}/)).toBeInTheDocument();
    });
  });
});
