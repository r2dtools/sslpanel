import { Button } from '@mui/base';
import Spinner from './Spinner';
import React from 'react';

type AppButtonProps = {
    loading?: boolean;
    fullWidth?: boolean;
    rounded?: boolean;
} & React.ComponentProps<typeof Button>

const AppButton = ({
    loading = false,
    fullWidth = false,
    rounded = false,
    children,
    ...props
}: AppButtonProps) => {
    let className = 'inline-flex items-center justify-center gap-2.5 bg-primary py-4 px-10 text-center font-medium text-white hover:bg-opacity-90 lg:px-8 xl:px-10';

    if (fullWidth) {
        className += ' w-full';
    }

    if (rounded) {
        className += ' rounded-md';
    }

    return (
        <Button
            className={className}
            {...props}
        >
            <Spinner show={loading} />
            {children}
        </Button>
    );
};

export default AppButton;
