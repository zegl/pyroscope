
export const mapNodeToPackageUrl = (node: ShamefulAny): string => {
    const name = node.name;
    if ( /^github.com/.test(name) || /^google.golang.org/.test(node.url) ) {
        const url = `https://pkg.go.dev/${name.split('.').slice(0,2).join('.')}`;
        console.log(url);
        return url
    }

    return "https://pkg.go.dev/" + name.split('.')[0]
}