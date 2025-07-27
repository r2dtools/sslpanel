import React, { useEffect, useState } from 'react';
import { Button, Checkbox, Drawer, HR, Label, Select, Spinner, Tabs, TextInput } from 'flowbite-react';
import { HiMiniLockClosed, HiMiniEnvelope } from 'react-icons/hi2';
import { AssignCertificatePayload, Domain, DomainSecurePayload } from '../types';
import { HTTP_CHALLENGE } from '../../server/constants';
import { StorageCertificateItem } from '../../certificate/types';

type SecureDomainProps = {
    guid: string,
    authToken: string,
    domain: Domain;
    open: boolean;
    issueLoading: boolean;
    assignLoading: boolean;
    certificates: StorageCertificateItem[];
    onIssueSubmit: (payload: DomainSecurePayload) => Promise<void>;
    onAssignSubmit: (payload: AssignCertificatePayload) => Promise<void>;
    onClose: () => void;
};

const SecureDomainDrawer: React.FC<SecureDomainProps> = ({
    guid,
    authToken,
    domain,
    open,
    issueLoading,
    assignLoading,
    certificates,
    onIssueSubmit,
    onAssignSubmit,
    onClose,
}) => {
    const certificateEmail = (domain.certificate?.emailaddresses || [])[0] || '';
    const [email, setEmail] = useState<string>(certificateEmail);
    const [assign, setAssign] = useState<boolean>(true);
    const [subjectsMap, setSubjectsMap] = useState<{ [key: string]: boolean }>({});

    useEffect(() => {
        setEmail(certificateEmail);
        setSubjectsMap({})
        setAssign(true);
    }, [domain, open]);

    const handleFormClose = (): void => {
        onClose();
    };

    const handleAliasChange = (alias: string, event: React.ChangeEvent<HTMLInputElement>): void => {
        setSubjectsMap({
            ...subjectsMap,
            [alias]: event.target.checked,
        });
    };

    const handleIssueSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const subjects = [];

        for (const [alias, checked] of Object.entries(subjectsMap)) {
            if (checked) {
                subjects.push(alias);
            }
        }

        await onIssueSubmit({
            guid,
            token: authToken,
            email,
            servername: domain.servername,
            subjects,
            webserver: domain.webserver,
            challengetype: HTTP_CHALLENGE,
            assign,
        });
        handleFormClose();
    };

    const selectOption = (item: StorageCertificateItem): string => `${item.certificate.cn} [${item.name}.pem, ${item.storage} storage]`;

    const handleAssignSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        // @ts-ignore
        const option = event.target.certificate.value;
        const entry = certificates.find((item: StorageCertificateItem) => option === selectOption(item));

        if (!entry) {
            return;
        }

        await onAssignSubmit({
            guid,
            token: authToken,
            servername: domain.servername,
            webserver: domain.webserver,
            name: entry.name,
            storage: entry.storage,
        });
        handleFormClose();
    }

    return (
        <Drawer className='z-[999] min-w-[450px]' open={open} onClose={handleFormClose} position='right'>
            <Drawer.Header title={`Secure ${domain.servername}`} titleIcon={HiMiniLockClosed} />
            <Drawer.Items>
                <Tabs aria-label="Tabs with underline" variant="underline">
                    <Tabs.Item active title="Issue certificate">
                        <p className='text-sm text-gray-500 dark:text-gray-400'>
                            <a href='https://letsencrypt.org' target='_blank' className='text-blue-600'>Let`s Encrypt</a> is a certificate authority (CA) that allows you to issue a free SSL/TLS certificate. By proceeding you acknowledge that you read and agree to the <a href='https://letsencrypt.org/repository' target='_blank'>Let`s Encrypt Terms of Service</a>.
                        </p>
                        <form onSubmit={handleIssueSubmit}>
                            <div className="mb-3 mt-3">
                                <Label htmlFor="email" className="mb-2 block">
                                    Email
                                </Label>
                                <TextInput
                                    id="email"
                                    name="email"
                                    placeholder="Email"
                                    required
                                    value={email}
                                    type='email'
                                    helperText={
                                        <>
                                            Email is used for registration and recovery contact.
                                        </>
                                    }
                                    rightIcon={HiMiniEnvelope}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => setEmail(event.target.value)}
                                />
                            </div>
                            <div className="flex max-w-md flex-col gap-3" id="checkbox">
                                <div className="flex items-center gap-2">
                                    <Checkbox id="servername" defaultChecked disabled />
                                    <Label htmlFor="servername" className="flex">
                                        {domain.servername}
                                    </Label>
                                </div>
                                {
                                    domain.aliases.map(alias => (
                                        <div className="flex items-center gap-2" key={alias}>
                                            <Checkbox
                                                id={alias}
                                                checked={subjectsMap[alias] || false}
                                                onChange={(event: React.ChangeEvent<HTMLInputElement>) => handleAliasChange(alias, event)}
                                                className='cursor-pointer'
                                            />
                                            <Label htmlFor={alias} className="flex cursor-pointer">
                                                {alias}
                                            </Label>
                                        </div>
                                    ))
                                }
                            </div>
                            <HR />
                            <div>
                                <div className="flex gap-2">
                                    <Checkbox id="assign" className='cursor-pointer' checked={assign} onChange={event => setAssign(event.target.checked)} />
                                    <div className='flex flex-col'>
                                        <Label htmlFor="assign" className="flex cursor-pointer">
                                            Assign certificate to the domain
                                        </Label>
                                        <div className="text-gray-500 dark:text-gray-300">
                                            <span className="text-xs font-normal">
                                                Clear checkbox if you are more conservative and want to make manual changes to the web server configuration. The certificate will be stored in storage.
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <Button className="w-full mt-4" color='blue' type='submit'>
                                {issueLoading ? <Spinner /> : 'Secure'}
                            </Button>
                        </form>
                    </Tabs.Item>
                    <Tabs.Item active title="Assign certificate">
                        <p className='text-sm text-gray-500 dark:text-gray-400'>
                            Assign certificate from the storage to the domain
                        </p>
                        <form onSubmit={handleAssignSubmit}>
                            <div className="mb-3 mt-3">
                                <Select id='certificate' className='cursor-pointer' required>
                                    {
                                        certificates.map((item: StorageCertificateItem) => {
                                            return <option className='cursor-pointer' key={`${item.name}__${item.storage}`}>{selectOption(item)}</option>;
                                        })
                                    }
                                </Select>
                            </div>
                            <Button className="w-full mt-4" color='blue' type='submit'>
                                {assignLoading ? <Spinner /> : 'Assign'}
                            </Button>
                        </form>
                    </Tabs.Item>
                </Tabs>
            </Drawer.Items>
        </Drawer>
    );
};

export default SecureDomainDrawer;
