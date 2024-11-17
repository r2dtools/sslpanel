import { Spinner } from 'flowbite-react';

const Loader = () => {
    return (
        <div className="absolute inset-0 flex items-center justify-center">
            <Spinner size='xl' color='success' />
        </div>
    );
};

export default Loader;
