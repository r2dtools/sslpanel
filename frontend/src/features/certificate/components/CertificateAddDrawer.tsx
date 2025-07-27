import { useEffect, useRef, useState } from 'react';
import { Badge, Button, Drawer, FileInput, Spinner, Tabs, TextInput } from 'flowbite-react';
import { HiMiniKey } from 'react-icons/hi2';
import isValidDomain from 'is-valid-domain';
import isValidFilename from 'valid-filename';
import { SelfSignedCertificateFormData } from '../types';

type CertificateAddProps = {
    open: boolean;
    uploading: boolean;
    generating: boolean;
    onClose: () => void;
    onUpload: (name: string, file: File) => Promise<any>;
    onGenerate: (data: SelfSignedCertificateFormData) => Promise<any>;
};

const CertificateAddDrawer: React.FC<CertificateAddProps> = ({
    open,
    uploading,
    generating,
    onClose,
    onUpload,
    onGenerate,
}) => {
    const [name, setName] = useState<string>('');
    const [commonName, setCommonName] = useState<string>('');
    const [email, setEmail] = useState<string>('');
    const [country, setCountry] = useState<string>('');
    const [province, setProvince] = useState<string>('');
    const [locality, setLocality] = useState<string>('');
    const [organization, setOrganization] = useState<string>('');
    const [alternativeNames, setAlternativeNames] = useState<string>('');

    const [nameErr, setNameErr] = useState<string>('');
    const [commonNameErr, setCommonNameErr] = useState<string>('');
    const [altNamesErr, setAltNamesErr] = useState<string>('');

    const [activeTab, setActiveTab] = useState<number>(0);

    let altNames = alternativeNames.split(',')
        .map((name: string) => name.trim())
        .filter((name: string) => !!name);
    altNames = [...(new Set(altNames))];

    const validAltNames = altNames.filter((name: string) => isValidDomain(name));

    const certUpload = useRef(null);

    useEffect(() => {
        setName('');
        setCommonName('');
        setEmail('');
        setCountry('');
        setProvince('');
        setLocality('');
        setOrganization('');
        setAlternativeNames('');
        resetCertUpload();
    }, [open, activeTab]);

    const resetCertUpload = (): void => {
        // @ts-ignore
        if (certUpload?.current?.value) {
            // @ts-ignore
            certUpload.current.value = '';
        }
    };

    const handleFormClose = (): void => {
        if (!generating) {
            resetCertUpload();
            onClose();
        }
    };

    const changeName = (name: string): void => {
        setName(name);
        setNameErr('');
    };

    const changeCommonName = (commonName: string): void => {
        setCommonName(commonName);
        setCommonNameErr('');
    };


    const changeAlternativeNames = (alternativeNames: string): void => {
        setAlternativeNames(alternativeNames);
        setAltNamesErr('');
    };

    const handleUploadSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        // @ts-ignore
        const certName = event.target.name.value;

        if (certName && !isValidFilename(certName)) {
            setNameErr('Invalid file name');

            return;
        }

        // @ts-ignore
        const file = event.target.certUpload.files[0] as File;
        const name = certName || file.name;

        await onUpload(name, file);
        handleFormClose();
    };

    const handleGenerateSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        // @ts-ignore
        const certName = event.target.name.value;

        if (certName && !isValidFilename(certName)) {
            setNameErr('Invalid file name');

            return;
        }

        // @ts-ignore
        const commonName = event.target.commonName.value;
        // @ts-ignore
        const email = event.target.email.value;
        // @ts-ignore
        const country = event.target.country.value;
        // @ts-ignore
        const province = event.target.province.value;
        // @ts-ignore
        const locality = event.target.locality.value;
        // @ts-ignore
        const organization = event.target.organization.value;

        if (!isValidDomain(commonName)) {
            setCommonNameErr('Invalid Common Name');

            return;
        }

        if (altNames.length !== validAltNames.length) {
            setAltNamesErr('Invalid Alternative Names');

            return;
        }

        await onGenerate({
            certName: certName,
            commonName: commonName,
            email: email || null,
            country: country || null,
            province: province || null,
            locality: locality || null,
            organization: organization || null,
            altNames: altNames,
        });
        handleFormClose();
    }

    return (
        <Drawer className='z-[999] min-w-[400px]' open={open} onClose={handleFormClose} position='right'>
            <Drawer.Header title='Add certificate' titleIcon={HiMiniKey} />
            <Drawer.Items>
                <Tabs aria-label="Tabs with underline" variant="underline" onActiveTabChange={tab => setActiveTab(tab)}>
                    <Tabs.Item active title="Upload PEM file">
                        <div className='mb-6 text-sm text-gray-500 dark:text-gray-400'>
                            <a href='https://en.wikipedia.org/wiki/Privacy-Enhanced_Mail' target='_blank' className='text-primary'>Privacy Enhanced Mail (PEM)</a> files are concatenated certificate containers frequently used in certificate installations when multiple certificates that form a complete chain are being imported as a single file.
                        </div>
                        <form onSubmit={handleUploadSubmit}>
                            <div className="mb-6 mt-3">
                                <TextInput
                                    id="name"
                                    name="name"
                                    placeholder="Certificate name"
                                    required
                                    value={name}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => changeName(event.target.value)}
                                    helperText={<>{nameErr}</>}
                                    color={nameErr ? 'failure' : undefined}
                                />
                            </div>
                            <div className='mb-6 mt-3'>
                                <FileInput
                                    id='cert-upload'
                                    name='certUpload'
                                    helperText='file with .pem extension'
                                    accept='.pem'
                                    ref={certUpload}
                                    required
                                />
                            </div>
                            <Button className="w-full" color='blue' type='submit'>
                                {uploading ? <Spinner /> : 'Upload'}
                            </Button>
                        </form>
                    </Tabs.Item>
                    <Tabs.Item title="Generate self-signed">
                        <div className='mb-6 text-sm text-gray-500 dark:text-gray-400'>
                            Self-signed SSL/TLS certificates can be used to set up temporary SSL servers. You can use it for test and development servers where security is not a big concern. Use the form below to generate a self-signed SSL/TLS certificate and key.
                        </div>
                        <form onSubmit={handleGenerateSubmit}>
                            <div className="mb-6 mt-3">
                                <TextInput
                                    id="name"
                                    name="name"
                                    placeholder="Certificate Name*"
                                    required
                                    value={name}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => setName(event.target.value)}
                                    helperText={<>{nameErr}</>}
                                    color={nameErr ? 'failure' : undefined}
                                />
                            </div>
                            <div className="mb-6 mt-3">
                                <TextInput
                                    id="commonName"
                                    name="commonName"
                                    placeholder="Common Name*"
                                    required
                                    value={commonName}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => changeCommonName(event.target.value)}
                                    helperText={commonNameErr ? <>{commonNameErr}</> : 'For example: example.com'}
                                    color={commonNameErr ? 'failure' : undefined}
                                />
                            </div>
                            <div className="mb-6 mt-3">
                                <TextInput
                                    id="email"
                                    name="email"
                                    placeholder="Email"
                                    type='email'
                                    value={email}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => setEmail(event.target.value)}
                                />
                            </div>
                            <div className="mb-6 mt-3">
                                <TextInput
                                    id="country"
                                    name="country"
                                    placeholder="Country"
                                    value={country}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => setCountry(event.target.value)}
                                />
                            </div>
                            <div className="mb-6 mt-3">
                                <TextInput
                                    id="province"
                                    name="province"
                                    placeholder="Province"
                                    value={province}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => setProvince(event.target.value)}
                                />
                            </div>
                            <div className="mb-6 mt-3">
                                <TextInput
                                    id="locality"
                                    name="locality"
                                    placeholder="Locality"
                                    value={locality}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => setLocality(event.target.value)}
                                />
                            </div>
                            <div className="mb-6 mt-3">
                                <TextInput
                                    id="organization"
                                    name="organization"
                                    placeholder="Organization"
                                    value={organization}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => setOrganization(event.target.value)}
                                />
                            </div>
                            <div className="mb-6 mt-3">
                                <TextInput
                                    id="alternativeNames"
                                    name="alternativeNames"
                                    placeholder="Alternative Names"
                                    value={alternativeNames}
                                    onChange={(event: React.ChangeEvent<HTMLInputElement>) => changeAlternativeNames(event.target.value)}
                                    helperText={altNamesErr ? <>{altNamesErr}</> : 'Comma separated alternative domain names. For example: www.example.com, webmail.example.com'}
                                    color={altNamesErr ? 'failure' : undefined}
                                />
                                {
                                    validAltNames.length ? (
                                        <div className='flex flex-row flex-wrap gap-2 mt-3'>
                                            {validAltNames.map((name: string) => <Badge key={name} color='gray'>{name}</Badge>)}
                                        </div>
                                    ) : null
                                }
                            </div>
                            <Button className="w-full" color='blue' type='submit'>
                                {generating ? <Spinner /> : 'Generate'}
                            </Button>
                        </form>
                    </Tabs.Item>
                </Tabs>

            </Drawer.Items>
        </Drawer >
    );
};

export default CertificateAddDrawer;
