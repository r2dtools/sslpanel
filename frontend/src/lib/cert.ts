export const isDnsNameSecure = (certificateDnsNames: string[], name: string): boolean =>
    Boolean(certificateDnsNames.find(certificateDnsName => certificateDnsName === name));
