import useLocalStorage, { SetValue } from '../../hooks/useLocalStorage';

const useAuthToken = (): [string | null, (value: SetValue<string | null>) => void] => {
    const [authToken, setAuthToken] = useLocalStorage<string | null>("r2panel-token", null);

    return [authToken, setAuthToken];
};

export default useAuthToken;
