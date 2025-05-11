import { useEffect } from 'react';
import useLocalStorage from './useLocalStorage';
import { ColorTheme } from '../features/account/types';

const useColorMode = () => {
  const [colorMode, setColorMode] = useLocalStorage('color-theme', ColorTheme.Dark);

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
