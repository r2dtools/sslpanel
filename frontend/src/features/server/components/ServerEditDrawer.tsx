import { useEffect, useState } from 'react';
import { Button, Drawer, Label, Spinner, TextInput } from 'flowbite-react';
import { HiMiniServer } from 'react-icons/hi2';
import { Server, ServerSavePayload } from '../types';

const defaultPort = 60150

type ServerEditProps = {
    server?: Server | null;
    authToken: string;
    open: boolean;
    loading: boolean;
    onSubmit: (server: ServerSavePayload) => void;
    onClose: () => void;
};

const ServerEditDrawer: React.FC<ServerEditProps> = ({
    server,
    authToken,
    open,
    loading,
    onSubmit,
    onClose,
}) => {
    const [name, setName] = useState<string>(server?.name || '');
    const [ipv4, setIpv4] = useState<string>(server?.ipv4_address || '');
    const [ipv6, setIpv6] = useState<string>(server?.ipv6_address || '');
    const [token, setToken] = useState<string>(server?.token || '');
    const [port, setPort] = useState<number>(server?.agent_port || defaultPort);
    const [ipError, setIpError] = useState<string>('');

    useEffect(() => {
        setName(open ? server?.name || '' : '');
        setIpv4(open ? server?.ipv4_address || '' : '');
        setIpv6(open ? server?.ipv6_address || '' : '');
        setPort(open ? server?.agent_port || defaultPort : defaultPort);
        setToken(open ? server?.token || '' : '');
    }, [server, open]);

    const handleFormClose = (): void => {
        onClose();
    };

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        if (ipv4 === '' && ipv6 === '') {
            setIpError('At leat one IP must be specified');

            return;
        }

        const saveServerPayload: ServerSavePayload = {
            authToken,
            name,
            ipv4_address: ipv4,
            ipv6_address: ipv6,
            agent_port: port,
            id: server?.id,
            token,
            guid: server?.guid,
        };

        onSubmit(saveServerPayload);
        handleFormClose();
    };

    return (
        <Drawer className='z-[999]' open={open} onClose={handleFormClose} position='right'>
            <Drawer.Header title={server?.id ? 'Edit server' : 'Add new server'} titleIcon={HiMiniServer} />
            <Drawer.Items>
                <form onSubmit={handleSubmit}>
                    <div className="mb-6 mt-3">
                        <Label htmlFor="name" className="mb-2 block">
                            Name
                        </Label>
                        <TextInput
                            id="name"
                            name="name"
                            placeholder="Server name"
                            required
                            value={name}
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => setName(event.target.value)}
                        />
                    </div>
                    <div className="mb-6">
                        <Label htmlFor="ipv4" className="mb-2 block">
                            IPv4
                        </Label>
                        <TextInput
                            id="ipv4"
                            name="ipv4"
                            placeholder="Server agent ipv4 address"
                            value={ipv4}
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => setIpv4(event.target.value)}
                            helperText={<>{ipError}</>}
                            color={ipError ? 'failure' : undefined}
                        />
                    </div>
                    <div className="mb-6">
                        <Label htmlFor="ipv6" className="mb-2 block">
                            IPv6
                        </Label>
                        <TextInput
                            id="ipv6"
                            name="ipv6"
                            placeholder="Server agent ipv6 address"
                            value={ipv6}
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => setIpv6(event.target.value)}
                            helperText={<>{ipError}</>}
                            color={ipError ? 'failure' : undefined}
                        />
                    </div>
                    <div className="mb-6">
                        <Label htmlFor="port" className="mb-2 block">
                            Port
                        </Label>
                        <TextInput
                            id="port"
                            name="port"
                            placeholder="Server agent port"
                            type='number'
                            required
                            value={port}
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => setPort(parseInt(event.target.value))}
                        />
                    </div>
                    <div className="mb-6 mt-3">
                        <Label htmlFor="token" className="mb-2 block">
                            Token
                        </Label>
                        <TextInput
                            id="token"
                            name="token"
                            placeholder="Server agent token"
                            required
                            value={token}
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => setToken(event.target.value)}
                        />
                    </div>
                    <Button className="w-full" color='blue' type='submit'>
                        {loading ? <Spinner /> : server ? 'Save' : 'Add'}
                    </Button>
                </form>
            </Drawer.Items>
        </Drawer>
    );
};

export default ServerEditDrawer;
