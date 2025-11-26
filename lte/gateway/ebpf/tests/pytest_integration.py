# pytest_integration.py
import pytest
import subprocess
from python.ebpf_manager import EbpfManager

@pytest.fixture
def ebpf_manager():
    return EbpfManager()

def test_manager_load(ebpf_manager):
    """Test that BPF manager can load programs"""
    result = ebpf_manager.load_programs()
    assert result is True

def test_session_sync(ebpf_manager):
    """Test GTP session sync"""
    result = ebpf_manager.sync_sessions()
    assert result is not None

def test_cleanup(ebpf_manager):
    """Test cleanup and detach"""
    result = ebpf_manager.cleanup()
    assert result is True
