#!/usr/bin/env python3

import json
import os

exceptions = {
    "proposals/README",
    "proposals/p004_fua-restrict-feature",
    "proposals/p006_subscriber_state_view",
    "proposals/p008_inbound_roaming_with_SubscriberDb",
    "proposals/p006_mandatory_integration_tests_for_each_PR.md",
    "proposals/p010_vendor_neutral_dp",
    "proposals/p011_intra_agw_mobility",
    "proposals/p012_resource-tagging",
    "proposals/sim_integration",
    "proposals/p021_mme_migrate_to_c++",
    "proposals/p022_enodebd_enhancements",
    "proposals/p023_magma_gtp_gateway",
    "proposals/p024_magma_settlement_service",
    "proposals/p025_magma_cdr_availability",
    "proposals/p026_magma_inbound_roaming_extensions",
    "proposals/qos_enforcement",
}


def get_implemented_sidebar_pages(version="latest"):
    version_prefix = _get_version_prefix(version)
    implemented_sidebar_pages = _extract_sidebar_pages(
        sidebar_json_path=_get_sidebar_json_path(version),
        version_prefix=version_prefix,
    )
    implemented_sidebar_pages = _remove_version_prefix(
        implemented_sidebar_pages=implemented_sidebar_pages,
        version_prefix=version_prefix,
    )
    return implemented_sidebar_pages


def get_all_pages(version="latest"):
    readme_path = _get_readme_path(version)
    all_pages = set()
    for root, dirs, filenames in os.walk(readme_path):
        for filename in filenames:
            if filename.endswith('.md'):
                doc_id = _extract_doc_id(filename, root, _get_version_prefix(version))
                all_pages.add(root.replace(f'{readme_path}/', '') + '/' + doc_id)
    return all_pages


def _extract_sidebar_pages(sidebar_json_path, version_prefix):
    implemented_sidebar_pages = set()
    with open(sidebar_json_path) as f:
        sidebars = json.load(f)[f'{version_prefix}docs']
        for k, v in sidebars.items():
            if isinstance(v[0], str):
                implemented_sidebar_pages = implemented_sidebar_pages.union(set(v))
            else:
                for item in v:
                    implemented_sidebar_pages = implemented_sidebar_pages.union(set(item['ids']))
    return implemented_sidebar_pages


def _remove_version_prefix(implemented_sidebar_pages, version_prefix):
    for page in implemented_sidebar_pages:
        if page.startswith(version_prefix):
            implemented_sidebar_pages.remove(page)
            implemented_sidebar_pages.add(page.replace(version_prefix, ''))
    return implemented_sidebar_pages


def _get_readme_path(version):
    if version == 'latest':
        readme_path = 'readmes'
    else:
        readme_path = f'docusaurus/versioned_docs/version-{version}'
    return readme_path


def _get_sidebar_json_path(version):
    if version == 'latest':
        sidebar_json_path = 'docusaurus/sidebars.json'
    else:
        sidebar_json_path = f'docusaurus/versioned_sidebars/version-{version}-sidebars.json'
    return sidebar_json_path


def _get_version_prefix(version):
    if version == 'latest':
        return ''
    else:
        return f'version-{version}-'


def _extract_doc_id(filename, root, version_prefix):
    path = os.path.join(root, filename)
    doc_id = ""
    with open(path) as f:
        lines = f.readlines()
        if len(lines) > 0 and lines[0].startswith('---'):
            for line in lines:
                if line.startswith('id: '):
                    doc_id = line.replace(f'id: {version_prefix}', '').rstrip('\n')
                    break
        else:
            doc_id = filename.replace('.md', '')
    return doc_id


if __name__ == '__main__':
    versions = ["latest", "1.8.0", "1.7.0"]
    pages_not_implemented = {v: set() for v in versions}
    for v in ["latest", "1.8.0", "1.7.0"]:
        all_pages = get_all_pages(version=v)
        sidebar_pages = get_implemented_sidebar_pages(version=v)
        for doc in sorted(all_pages):
            if doc not in sidebar_pages.union(exceptions):
                pages_not_implemented[v].add(doc)

    sidebars_missing = False
    for v in versions:
        if len(pages_not_implemented[v]) > 0:
            sidebars_missing = True
            print(f"Missing pages for {v}: {pages_not_implemented[v]}")
    if sidebars_missing:
        exit(1)
