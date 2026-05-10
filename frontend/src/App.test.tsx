import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import App from './App';
import { ThemeProvider } from './context/ThemeContext';
import { AppWrapper } from './components/common/PageMeta';

describe('App', () => {
  it('renders without crashing', () => {
    render(
      <ThemeProvider>
        <AppWrapper>
          <App />
        </AppWrapper>
      </ThemeProvider>
    );
    expect(document.body).toBeDefined();
  });
});
