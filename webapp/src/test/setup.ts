import '@testing-library/jest-dom';
import { vi } from 'vitest';

// Mock window.registerPlugin
global.window = global.window || {};
global.window.registerPlugin = vi.fn();

// Mock Mattermost global objects
global.window.React = require('react');
global.window.ReactDOM = require('react-dom');
