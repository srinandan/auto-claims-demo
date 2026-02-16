
import asyncio
import os
import sys

# Ensure app can be imported
sys.path.append(os.getcwd())

from app.agent_config import agent_card as file_agent_card

async def verify():
    print("Verifying Agent Card Skills from agent_config.py...")
    
    skills = file_agent_card.skills
    found = False
    for skill in skills:
        if skill.id == "assess_severity":
            print(f"Found skill: {skill.name}")
            print(f"Description: {skill.description}")
            print(f"Examples: {skill.examples}")
            found = True
            break
    
    if found:
        print("Verification SUCCESS: Agent Skill found.")
    else:
        print("Verification FAILED: Agent Skill NOT found.")
        exit(1)

if __name__ == "__main__":
    asyncio.run(verify())
