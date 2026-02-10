import { describe, it, expect, vi, afterEach } from 'vitest';
import { render, screen, fireEvent, cleanup } from '@testing-library/react';
import TTLMenu from './ttl_menu';

describe('TTLMenu', () => {
  afterEach(() => {
    cleanup();
  });

  const defaultProps = {
    durations: [
      { label: '5 minutes', value: '5m' },
      { label: '15 minutes', value: '15m' },
      { label: '1 hour', value: '1h' },
      { label: '1 day', value: '1d' },
    ],
    selected: null as string | null,
    onSelect: vi.fn(),
    onClear: vi.fn(),
  };

  describe('rendering', () => {
    it('should render menu header', () => {
      render(<TTLMenu {...defaultProps} />);

      expect(screen.getByText('Message expires in...')).toBeInTheDocument();
    });

    it('should render all duration options', () => {
      render(<TTLMenu {...defaultProps} />);

      expect(screen.getByText('5 minutes')).toBeInTheDocument();
      expect(screen.getByText('15 minutes')).toBeInTheDocument();
      expect(screen.getByText('1 hour')).toBeInTheDocument();
      expect(screen.getByText('1 day')).toBeInTheDocument();
    });

    it('should render correct container class', () => {
      const { container } = render(<TTLMenu {...defaultProps} />);

      expect(container.querySelector('.ttl-menu')).toBeInTheDocument();
    });
  });

  describe('selection display', () => {
    it('should show checkmark for selected duration', () => {
      const props = { ...defaultProps, selected: '5m' };
      render(<TTLMenu {...props} />);

      const fiveMinOption = screen.getByText('5 minutes').closest('.ttl-option');
      expect(fiveMinOption).toHaveClass('selected');
      expect(fiveMinOption?.querySelector('.ttl-option-check')).toHaveTextContent('✓');
    });

    it('should not show checkmark for unselected durations', () => {
      const props = { ...defaultProps, selected: '5m' };
      render(<TTLMenu {...props} />);

      const fifteenMinOption = screen.getByText('15 minutes').closest('.ttl-option');
      expect(fifteenMinOption).not.toHaveClass('selected');
      expect(fifteenMinOption?.querySelector('.ttl-option-check')).toBeNull();
    });

    it('should not show any checkmarks when nothing is selected', () => {
      render(<TTLMenu {...defaultProps} />);

      expect(screen.queryByText('✓')).toBeNull();
    });

    it('should apply selected class to selected option', () => {
      const props = { ...defaultProps, selected: '1h' };
      render(<TTLMenu {...props} />);

      const oneHourOption = screen.getByText('1 hour').closest('.ttl-option');
      expect(oneHourOption).toHaveClass('selected');
    });
  });

  describe('option clicking', () => {
    it('should call onSelect with correct value when option is clicked', () => {
      render(<TTLMenu {...defaultProps} />);

      const fiveMinOption = screen.getByText('5 minutes').closest('.ttl-option');
      fireEvent.click(fiveMinOption!);

      expect(defaultProps.onSelect).toHaveBeenCalledWith('5m');
    });

    it('should call onSelect for 15 minutes option', () => {
      render(<TTLMenu {...defaultProps} />);

      const fifteenMinOption = screen.getByText('15 minutes').closest('.ttl-option');
      fireEvent.click(fifteenMinOption!);

      expect(defaultProps.onSelect).toHaveBeenCalledWith('15m');
    });

    it('should call onSelect for 1 hour option', () => {
      render(<TTLMenu {...defaultProps} />);

      const oneHourOption = screen.getByText('1 hour').closest('.ttl-option');
      fireEvent.click(oneHourOption!);

      expect(defaultProps.onSelect).toHaveBeenCalledWith('1h');
    });

    it('should call onSelect for 1 day option', () => {
      render(<TTLMenu {...defaultProps} />);

      const oneDayOption = screen.getByText('1 day').closest('.ttl-option');
      fireEvent.click(oneDayOption!);

      expect(defaultProps.onSelect).toHaveBeenCalledWith('1d');
    });

    it('should not call onClear when clicking a duration option', () => {
      render(<TTLMenu {...defaultProps} />);

      const fiveMinOption = screen.getByText('5 minutes').closest('.ttl-option');
      fireEvent.click(fiveMinOption!);

      expect(defaultProps.onClear).not.toHaveBeenCalled();
    });
  });

  describe('clear option', () => {
    it('should not show clear option when nothing is selected', () => {
      render(<TTLMenu {...defaultProps} />);

      expect(screen.queryByText('Disable TTL')).toBeNull();
    });

    it('should show clear option when a duration is selected', () => {
      const props = { ...defaultProps, selected: '5m' };
      render(<TTLMenu {...props} />);

      expect(screen.getByText('Disable TTL')).toBeInTheDocument();
    });

    it('should apply clear option class', () => {
      const props = { ...defaultProps, selected: '5m' };
      render(<TTLMenu {...props} />);

      const clearOption = screen.getByText('Disable TTL').closest('.ttl-option');
      expect(clearOption).toHaveClass('ttl-option-clear');
    });

    it('should call onClear when clear option is clicked', () => {
      const props = { ...defaultProps, selected: '5m' };
      render(<TTLMenu {...props} />);

      const clearOption = screen.getByText('Disable TTL').closest('.ttl-option');
      fireEvent.click(clearOption!);

      expect(props.onClear).toHaveBeenCalled();
      expect(props.onSelect).not.toHaveBeenCalled();
    });

    it('should not call onSelect when clear option is clicked', () => {
      const props = { ...defaultProps, selected: '5m' };
      render(<TTLMenu {...props} />);

      const clearOption = screen.getByText('Disable TTL').closest('.ttl-option');
      fireEvent.click(clearOption!);

      expect(props.onSelect).not.toHaveBeenCalled();
    });
  });

  describe('edge cases', () => {
    it('should handle empty durations array', () => {
      const props = { ...defaultProps, durations: [] };
      const { container } = render(<TTLMenu {...props} />);

      expect(container.querySelector('.ttl-menu-options')).toBeInTheDocument();
      expect(container.querySelector('.ttl-option')).toBeNull();
    });

    it('should handle single duration', () => {
      const props = {
        ...defaultProps,
        durations: [{ label: '5 minutes', value: '5m' }],
      };
      render(<TTLMenu {...props} />);

      expect(screen.getByText('5 minutes')).toBeInTheDocument();
    });

    it('should handle many durations', () => {
      const manyDurations = [
        { label: '1 minute', value: '1m' },
        { label: '5 minutes', value: '5m' },
        { label: '10 minutes', value: '10m' },
        { label: '15 minutes', value: '15m' },
        { label: '30 minutes', value: '30m' },
        { label: '1 hour', value: '1h' },
        { label: '2 hours', value: '2h' },
        { label: '6 hours', value: '6h' },
        { label: '12 hours', value: '12h' },
        { label: '1 day', value: '1d' },
      ];
      const props = { ...defaultProps, durations: manyDurations };
      render(<TTLMenu {...props} />);

      expect(screen.getByText('1 minute')).toBeInTheDocument();
      expect(screen.getByText('12 hours')).toBeInTheDocument();
      expect(screen.getByText('1 day')).toBeInTheDocument();
    });
  });
});
