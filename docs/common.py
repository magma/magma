import os


def get_readme_path(version):
    """Get the path to the readme folder for a given version."""
    if version == 'latest':
        return 'readmes'
    return f'docusaurus/versioned_docs/version-{version}'


def extract_doc_id(filename, root, version_prefix):
    """Extract the doc id from a given file."""
    path = os.path.join(root, filename)
    doc_id = ""
    with open(path) as f:
        lines = f.readlines()
        if lines and lines[0].startswith('---'):
            for line in lines:
                if line.startswith('id: '):
                    doc_id = line.replace(f'id: {version_prefix}', '').rstrip('\n')
                    break
        else:
            doc_id = filename.replace('.md', '')
    return doc_id


def get_version_prefix(version):
    """Get the version prefix for a given version."""
    if version == 'latest':
        return ''
    return f'version-{version}-'
