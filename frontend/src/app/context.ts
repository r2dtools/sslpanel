import { createContext } from 'react';
import { RouteItem } from './types';

const RoutesContext = createContext<RouteItem[]>([]);

export default RoutesContext;
