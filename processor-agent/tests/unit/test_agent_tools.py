try:
    from app.app_utils.repair_costs import generate_repair_cost
except (ImportError, ModuleNotFoundError):
    from app_utils.repair_costs import generate_repair_cost

def test_generate_repair_cost_simple():
    """Test with 'simple' severity."""
    result = generate_repair_cost("simple")
    assert result["total_cost"] == 700.00
    assert result["total_labor"] == 200.00
    assert result["total_parts"] == 500.00
    assert len(result["items"]) == 3
    assert result["items"][0]["part"] == "Bumper Repair"

def test_generate_repair_cost_complex():
    """Test with 'complex' severity."""
    result = generate_repair_cost("complex")
    assert result["total_cost"] == 4500.00
    assert result["total_labor"] == 1000.00
    assert result["total_parts"] == 3500.00
    assert len(result["items"]) == 4
    assert result["items"][0]["part"] == "Fender Replacement"

def test_generate_repair_cost_case_insensitive():
    """Test that it handles case-insensitive input."""
    result = generate_repair_cost("SIMPLE")
    assert result["total_cost"] == 700.00

    result = generate_repair_cost("CoMpLeX")
    assert result["total_cost"] == 4500.00

def test_generate_repair_cost_substring():
    """Test that it handles strings containing 'simple'."""
    result = generate_repair_cost("This is a simple repair")
    assert result["total_cost"] == 700.00

def test_generate_repair_cost_default_to_complex():
    """Test that it defaults to complex if 'simple' is not present."""
    result = generate_repair_cost("moderate")
    assert result["total_cost"] == 4500.00

    result = generate_repair_cost("")
    assert result["total_cost"] == 4500.00
