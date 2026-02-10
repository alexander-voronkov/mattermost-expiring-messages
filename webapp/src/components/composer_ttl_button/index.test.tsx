import { describe, it, expect, vi, afterEach } from 'vitest';
import { render, screen, fireEvent, cleanup } from '@testing-library/react';
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import ComposerTTLButton from './index';

describe('ComposerTTLButton', () => {
  afterEach(() => {
    cleanup();
    delete (window as any).setSelectedTTLDuration;
  });

  const createMockStore = () => {
    return configureStore({
      reducer: {
        plugins: () => ({
          'com.fambear.expiring-messages': {
            selectedDuration: null,
          },
        }),
      },
    });
  };

  const renderComponent = (store = createMockStore()) => {
    return render(
      <Provider store={store}>
        <ComposerTTLButton store={store} />
      </Provider>
    );
  };

  describe('rendering', () => {
    it('should render flame icon button', () => {
      renderComponent();

      expect(screen.getByText('🔥')).toBeInTheDocument();
    });

    it('should render button with correct class', () => {
      const { container } = renderComponent();

      expect(container.querySelector('.composer-ttl-button')).toBeInTheDocument();
    });

    it('should not have active class when no duration selected', () => {
      const { container } = renderComponent();

      const button = container.querySelector('.composer-ttl-button');
      expect(button).not.toHaveClass('active');
    });

    it('should have active class when duration is selected', () => {
      const store = configureStore({
        reducer: {
          plugins: () => ({
            'com.fambear.expiring-messages': {
              selectedDuration: '5m',
            },
          }),
        },
      });

      const { container } = render(
        <Provider store={store}>
          <ComposerTTLButton store={store} />
        </Provider>
      );

      const button = container.querySelector('.composer-ttl-button');
      expect(button).toHaveClass('active');
    });

    it('should show default title when no duration selected', () => {
      const { container } = renderComponent();

      const button = container.querySelector('.composer-ttl-button');
      expect(button).toHaveAttribute('title', 'Set message expiration');
    });

    it('should show duration in title when selected', () => {
      const store = configureStore({
        reducer: {
          plugins: () => ({
            'com.fambear.expiring-messages': {
              selectedDuration: '5m',
            },
          }),
        },
      });

      const { container } = render(
        <Provider store={store}>
          <ComposerTTLButton store={store} />
        </Provider>
      );

      const button = container.querySelector('.composer-ttl-button');
      expect(button).toHaveAttribute('title', 'TTL: 5 minutes');
    });
  });

  describe('menu interactions', () => {
    it('should not show menu initially', () => {
      const { container } = renderComponent();

      expect(container.querySelector('.ttl-menu')).toBeNull();
      expect(container.querySelector('.ttl-menu-overlay')).toBeNull();
    });

    it('should show menu when button is clicked', () => {
      const { container } = renderComponent();

      const button = container.querySelector('.composer-ttl-button')!;
      fireEvent.click(button);

      expect(container.querySelector('.ttl-menu')).toBeInTheDocument();
      expect(container.querySelector('.ttl-menu-overlay')).toBeInTheDocument();
    });

    it('should hide menu when button is clicked again', () => {
      const { container } = renderComponent();

      const button = container.querySelector('.composer-ttl-button')!;
      fireEvent.click(button);
      fireEvent.click(button);

      expect(container.querySelector('.ttl-menu')).toBeNull();
    });

    it('should hide menu when overlay is clicked', () => {
      const { container } = renderComponent();

      const button = container.querySelector('.composer-ttl-button')!;
      fireEvent.click(button);

      const overlay = container.querySelector('.ttl-menu-overlay')!;
      fireEvent.click(overlay);

      expect(container.querySelector('.ttl-menu')).toBeNull();
    });
  });

  describe('duration selection', () => {
    it('should set window.setSelectedTTLDuration when option is clicked', () => {
      const { container } = renderComponent();

      const button = container.querySelector('.composer-ttl-button')!;
      fireEvent.click(button);

      const fiveMinOption = screen.getByText('5 minutes').closest('.ttl-option');
      fireEvent.click(fiveMinOption!);

      expect((window as any).setSelectedTTLDuration).toBe('5m');
    });

    it('should show all predefined duration options', () => {
      const { container } = renderComponent();

      const button = container.querySelector('.composer-ttl-button')!;
      fireEvent.click(button);

      expect(screen.getByText('5 minutes')).toBeInTheDocument();
      expect(screen.getByText('15 minutes')).toBeInTheDocument();
      expect(screen.getByText('1 hour')).toBeInTheDocument();
      expect(screen.getByText('1 day')).toBeInTheDocument();
    });

    it('should highlight selected option in menu', () => {
      const store = configureStore({
        reducer: {
          plugins: () => ({
            'com.fambear.expiring-messages': {
              selectedDuration: '15m',
            },
          }),
        },
      });

      const { container } = render(
        <Provider store={store}>
          <ComposerTTLButton store={store} />
        </Provider>
      );

      const button = container.querySelector('.composer-ttl-button')!;
      fireEvent.click(button);

      const fifteenMinOption = screen.getByText('15 minutes').closest('.ttl-option');
      expect(fifteenMinOption).toHaveClass('selected');
    });
  });

  describe('clearing selection', () => {
    it('should show clear option when duration is selected', () => {
      const store = configureStore({
        reducer: {
          plugins: () => ({
            'com.fambear.expiring-messages': {
              selectedDuration: '5m',
            },
          }),
        },
      });

      const { container } = render(
        <Provider store={store}>
          <ComposerTTLButton store={store} />
        </Provider>
      );

      const button = container.querySelector('.composer-ttl-button')!;
      fireEvent.click(button);

      expect(screen.getByText('Disable TTL')).toBeInTheDocument();
    });

    it('should remove window.setSelectedTTLDuration when clear is clicked', () => {
      const store = configureStore({
        reducer: {
          plugins: () => ({
            'com.fambear.expiring-messages': {
              selectedDuration: '5m',
            },
          }),
        },
      });

      (window as any).setSelectedTTLDuration = '5m';

      const { container } = render(
        <Provider store={store}>
          <ComposerTTLButton store={store} />
        </Provider>
      );

      const button = container.querySelector('.composer-ttl-button')!;
      fireEvent.click(button);

      const clearOption = screen.getByText('Disable TTL').closest('.ttl-option');
      fireEvent.click(clearOption!);

      expect((window as any).setSelectedTTLDuration).toBeUndefined();
    });

    it('should close menu after clearing selection', () => {
      const store = configureStore({
        reducer: {
          plugins: () => ({
            'com.fambear.expiring-messages': {
              selectedDuration: '5m',
            },
          }),
        },
      });

      const { container } = render(
        <Provider store={store}>
          <ComposerTTLButton store={store} />
        </Provider>
      );

      const button = container.querySelector('.composer-ttl-button')!;
      fireEvent.click(button);

      const clearOption = screen.getByText('Disable TTL').closest('.ttl-option');
      fireEvent.click(clearOption!);

      expect(container.querySelector('.ttl-menu')).toBeNull();
    });
  });

  describe('edge cases', () => {
    it('should handle store without plugin state', () => {
      const store = configureStore({
        reducer: {
          plugins: () => ({}),
        },
      });

      const { container } = render(
        <Provider store={store}>
          <ComposerTTLButton store={store} />
        </Provider>
      );

      expect(container.querySelector('.composer-ttl-button')).toBeInTheDocument();
    });

    it('should handle store with empty plugin state', () => {
      const store = configureStore({
        reducer: {
          plugins: () => null,
        },
      });

      const { container } = render(
        <Provider store={store}>
          <ComposerTTLButton store={store} />
        </Provider>
      );

      expect(container.querySelector('.composer-ttl-button')).toBeInTheDocument();
    });

    it('should handle unrecognized duration value', () => {
      const store = configureStore({
        reducer: {
          plugins: () => ({
            'com.fambear.expiring-messages': {
              selectedDuration: '99m',
            },
          }),
        },
      });

      const { container } = render(
        <Provider store={store}>
          <ComposerTTLButton store={store} />
        </Provider>
      );

      const button = container.querySelector('.composer-ttl-button');
      expect(button).toHaveAttribute('title', 'TTL: 99m');
    });
  });
});
