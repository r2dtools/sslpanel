import { useEffect } from 'react';
import useLocalStorage from './useLocalStorage';
import { ColorTheme } from '../types/theme';

const useColorMode = () => {
  const [colorMode, setColorMode] = useLocalStorage('color-theme', ColorTheme.Light);

  useEffect(() => {
    const className = ColorTheme.Dark;
    const bodyClass = window.document.body.classList;

    colorMode === ColorTheme.Dark
      ? bodyClass.add(className)
      : bodyClass.remove(className);
  }, [colorMode]);

  return [colorMode, setColorMode];
};

export default useColorMode;
