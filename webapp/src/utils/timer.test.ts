import { describe, it, expect } from 'vitest';
import { formatTimeRemaining, parseDuration } from './timer';

describe('formatTimeRemaining', () => {
  it('should format seconds correctly', () => {
    expect(formatTimeRemaining(5000)).toBe('00:05');
    expect(formatTimeRemaining(30000)).toBe('00:30');
    expect(formatTimeRemaining(59000)).toBe('00:59');
  });

  it('should format minutes and seconds correctly', () => {
    expect(formatTimeRemaining(60000)).toBe('01:00');
    expect(formatTimeRemaining(60000 + 5000)).toBe('01:05');
    expect(formatTimeRemaining(5 * 60000 + 30000)).toBe('05:30');
    expect(formatTimeRemaining(59 * 60000 + 59000)).toBe('59:59');
  });

  it('should format hours, minutes, and seconds correctly', () => {
    expect(formatTimeRemaining(3600000)).toBe('1:00:00');
    expect(formatTimeRemaining(3600000 + 60000 + 5000)).toBe('1:01:05');
    expect(formatTimeRemaining(2 * 3600000 + 30 * 60000 + 15000)).toBe('2:30:15');
  });

  it('should handle zero milliseconds', () => {
    expect(formatTimeRemaining(0)).toBe('00:00');
  });

  it('should handle sub-second values', () => {
    expect(formatTimeRemaining(500)).toBe('00:00');
    expect(formatTimeRemaining(999)).toBe('00:00');
  });

  it('should handle large values', () => {
    expect(formatTimeRemaining(24 * 3600000)).toBe('24:00:00');
    expect(formatTimeRemaining(100 * 3600000 + 15 * 60000 + 30000)).toBe('100:15:30');
  });

  it('should pad single digit numbers', () => {
    expect(formatTimeRemaining(60000 + 5000)).toBe('01:05');
    expect(formatTimeRemaining(3600000 + 60000 + 5000)).toBe('1:01:05');
  });
});

describe('parseDuration', () => {
  describe('valid durations', () => {
    it('should parse minutes', () => {
      expect(parseDuration('5m')).toBe(5 * 60 * 1000);
      expect(parseDuration('15m')).toBe(15 * 60 * 1000);
      expect(parseDuration('60m')).toBe(60 * 60 * 1000);
    });

    it('should parse hours', () => {
      expect(parseDuration('1h')).toBe(60 * 60 * 1000);
      expect(parseDuration('2h')).toBe(2 * 60 * 60 * 1000);
      expect(parseDuration('24h')).toBe(24 * 60 * 60 * 1000);
    });

    it('should parse days', () => {
      expect(parseDuration('1d')).toBe(24 * 60 * 60 * 1000);
      expect(parseDuration('3d')).toBe(3 * 24 * 60 * 60 * 1000);
      expect(parseDuration('7d')).toBe(7 * 24 * 60 * 60 * 1000);
    });

    it('should handle large values', () => {
      expect(parseDuration('999m')).toBe(999 * 60 * 1000);
      expect(parseDuration('999h')).toBe(999 * 60 * 60 * 1000);
      expect(parseDuration('999d')).toBe(999 * 24 * 60 * 60 * 1000);
    });
  });

  describe('invalid durations', () => {
    it('should return 0 for invalid format', () => {
      expect(parseDuration('')).toBe(0);
      expect(parseDuration('5')).toBe(0);
      expect(parseDuration('m')).toBe(0);
      expect(parseDuration('5x')).toBe(0);
      expect(parseDuration('invalid')).toBe(0);
      expect(parseDuration('-5m')).toBe(0);
      expect(parseDuration('5.5m')).toBe(0);
      expect(parseDuration('5  m')).toBe(0);
    });

    it('should return 0 for invalid units', () => {
      expect(parseDuration('5s')).toBe(0);
      expect(parseDuration('5w')).toBe(0);
      expect(parseDuration('5M')).toBe(0);
      expect(parseDuration('5H')).toBe(0);
    });

    it('should return 0 for non-numeric values', () => {
      expect(parseDuration('am')).toBe(0);
      expect(parseDuration('abc')).toBe(0);
      expect(parseDuration('1m2h')).toBe(0);
    });
  });

  describe('edge cases', () => {
    it('should handle 1m', () => {
      expect(parseDuration('1m')).toBe(60 * 1000);
    });

    it('should handle 1h', () => {
      expect(parseDuration('1h')).toBe(60 * 60 * 1000);
    });

    it('should handle 1d', () => {
      expect(parseDuration('1d')).toBe(24 * 60 * 60 * 1000);
    });

    it('should handle 0m (edge case)', () => {
      expect(parseDuration('0m')).toBe(0);
    });
  });
});
